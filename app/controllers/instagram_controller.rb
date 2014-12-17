class InstagramController < ApplicationController
  before_action :require_signed_in

  def connect
    redirect_to Instagram.authorize_url(redirect_uri: ig_callback_url)
  end

  def callback
    response = Instagram.get_access_token(params[:code], redirect_uri: ig_callback_url)
    if InstagramConnection.create(user_id: current_user.id, access_token: response.access_token)
      flash[:success] = 'Instagram successfully connected'
      redirect_to root_path
    end
  end
end
