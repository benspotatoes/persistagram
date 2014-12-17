class InstagramConnection < ActiveRecord::Base
  LIKED_PHOTOS_ENDPOINT = 'https://api.instagram.com/v1/users/self/media/liked'

  def self.liked_photos_endpoint(access_token:)
    "#{LIKED_PHOTOS_ENDPOINT}?access_token=#{access_token}"
  end
end
