class ApplicationController < ActionController::Base
  include Clearance::Controller
  # Prevent CSRF attacks by raising an exception.
  # For APIs, you may want to use :null_session instead.
  protect_from_forgery with: :exception

  def require_signed_in
    unless signed_in?
      flash[:error] = 'You must be signed in'
      redirect_to sign_in_path
    end
  end
end
