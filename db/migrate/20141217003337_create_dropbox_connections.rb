class CreateDropboxConnections < ActiveRecord::Migration
  def change
    create_table :dropbox_connections do |t|
      t.integer :user_id, null: false, unique: true
      t.string :access_token, null: false, unique: true

      t.timestamps
    end

    add_index :dropbox_connections, [:user_id], unique: true
    add_index :dropbox_connections, [:access_token], unique: true
    add_index :dropbox_connections, [:user_id, :access_token], unique: true
  end
end
