class CreateInstagramConnections < ActiveRecord::Migration
  def change
    create_table :instagram_connections do |t|

      t.timestamps
    end
  end
end
