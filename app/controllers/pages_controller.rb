class PagesController < ApplicationController
  before_action :require_signed_in

  def home
    if signed_in?
      @ig_conn = InstagramConnection.find_by(user_id: current_user.id)
      @db_conn = DropboxConnection.find_by(user_id: current_user.id)
    end
  end

  def import
  end
end
