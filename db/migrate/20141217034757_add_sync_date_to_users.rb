class AddSyncDateToUsers < ActiveRecord::Migration
  def change
    add_column :users, :last_sync, :datetime
  end
end
