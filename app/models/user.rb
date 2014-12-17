class User < ActiveRecord::Base
  include Clearance::User

  SYNC_LIMIT_DAYS = 7

  def can_sync?
    since_last = ((Time.now - last_sync) / 86400).round
    return since_last > SYNC_LIMIT_DAYS.days, SYNC_LIMIT_DAYS - since_last
  end
end
