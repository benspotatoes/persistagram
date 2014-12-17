class DropboxController < ApplicationController
  before_action :require_signed_in

  def connect
    redirect_to db_auth.start
  end

  def callback
    access_token, user_id, url_state = db_auth.finish(params)
    if DropboxConnection.create(user_id: current_user.id, access_token: access_token)
      flash[:success] = 'Dropbox successfully connected'
      redirect_to root_path
    end
  end

  private

  def db_auth
    require 'dropbox_sdk' if !defined?(DropboxOAuth2Flow)
    DropboxOAuth2Flow.new(
      Rails.application.secrets.dropbox_client_id,
      Rails.application.secrets.dropbox_client_secret,
      db_callback_url,
      session,
      :dropbox_auth_csrf_token)
  end
end
