class CreateUserMetrics < ActiveRecord::Migration
  def change
    create_table :user_metrics do |t|
      t.integer :user_id, unique: false, null: false
      t.integer :files_saved, unique: false, null: false

      t.timestamps
    end

    add_index :user_metrics, [:user_id]
  end
end
