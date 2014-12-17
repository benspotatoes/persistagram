class CreateDropboxConnections < ActiveRecord::Migration
  def change
    create_table :dropbox_connections do |t|

      t.timestamps
    end
  end
end
