require 'open-uri'

class PagesController < ApplicationController
  before_action :require_signed_in
  before_action :set_ig_conn
  before_action :set_db_conn

  def home
  end

  def sync
    if [@ig_conn, @db_conn, @db].any? { |obj| obj.nil? }
      flash[:error] = 'An error occurred'
      redirect_to root_path
      return
    end

    should_proceed, return_in = current_user.can_sync?
    if !should_proceed
      flash[:error] = "Cannot sync, try again in #{return_in} days"
      redirect_to root_path
      return
    end

    metrics = {req_count: 1, files_saved: 0}
    likes = {}
    url = InstagramConnection.liked_photos_endpoint(access_token: @ig_conn.access_token)

    Rails.logger.info("Getting likes for user id #{current_user.id}")
    while !url.nil? do
      Rails.logger.debug("Iteration ##{metrics[:req_count]}")
      resp = JSON.parse(Faraday.new(url: url).get.body)

      url = resp['pagination']['next_url']

      resp['data'].each do |like|
        username = like['user']['username']
        likes[username] ||= []

        videos = like['videos']
        images = like['images']

        if videos
          likes[username] << videos['standard_resolution']['url']
        elsif images
          likes[username] << images['standard_resolution']['url']
        else
          Rails.logger.error('Invalid media')
        end
      end
      metrics[:req_count] += 1
    end

    Rails.logger.info("Saving liked media for user id #{current_user.id}")
    likes.each do |username, media|
      media.each do |item|
        item.match(/.*\/(\w*)\.(\w{3})/)
        filename = $1
        extension = $2
        path = "/iglikes-#{current_user.id}/#{username}/#{filename}.#{extension}"

        if @db.search("/iglikes-#{current_user.id}/#{username}", "#{filename}.#{extension}").empty?
          Rails.logger.debug("Created file: #{username} - #{filename}.#{extension}")
          @db.put_file(path, open(item).read)
          metrics[:files_saved] += 1
        else
          Rails.logger.debug("File exists: #{username} - #{filename}.#{extension}")
        end
      end
    end

    current_user.user_metrics.create!(files_saved: metrics[:files_saved])
    Rails.logger.info("#{metrics[:files_saved]} files saved")

    current_user.last_sync = Time.now
    current_user.save!
  end

  private

  def set_ig_conn
    @ig_conn = InstagramConnection.find_by(user_id: current_user.try(:id))
  end

  def set_db_conn
    require 'dropbox_sdk' if !defined?(DropboxClient)
    @db_conn = DropboxConnection.find_by(user_id: current_user.try(:id))
    @db = DropboxClient.new(@db_conn.access_token) if @db_conn
  end
end
