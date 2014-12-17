# encoding: UTF-8
# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20141217034757) do

  create_table "dropbox_connections", force: true do |t|
    t.integer  "user_id",      null: false
    t.string   "access_token", null: false
    t.datetime "created_at"
    t.datetime "updated_at"
  end

  add_index "dropbox_connections", ["access_token"], name: "index_dropbox_connections_on_access_token", unique: true
  add_index "dropbox_connections", ["user_id", "access_token"], name: "index_dropbox_connections_on_user_id_and_access_token", unique: true
  add_index "dropbox_connections", ["user_id"], name: "index_dropbox_connections_on_user_id", unique: true

  create_table "instagram_connections", force: true do |t|
    t.integer  "user_id",      null: false
    t.string   "access_token", null: false
    t.datetime "created_at"
    t.datetime "updated_at"
  end

  add_index "instagram_connections", ["access_token"], name: "index_instagram_connections_on_access_token", unique: true
  add_index "instagram_connections", ["user_id", "access_token"], name: "index_instagram_connections_on_user_id_and_access_token", unique: true
  add_index "instagram_connections", ["user_id"], name: "index_instagram_connections_on_user_id", unique: true

  create_table "users", force: true do |t|
    t.datetime "created_at",                     null: false
    t.datetime "updated_at",                     null: false
    t.string   "email",                          null: false
    t.string   "encrypted_password", limit: 128, null: false
    t.string   "confirmation_token", limit: 128
    t.string   "remember_token",     limit: 128, null: false
    t.datetime "last_sync"
  end

  add_index "users", ["email"], name: "index_users_on_email"
  add_index "users", ["remember_token"], name: "index_users_on_remember_token"

end
