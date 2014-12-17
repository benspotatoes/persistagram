Rails.application.routes.draw do
  root 'pages#home'

  scope 'ig' do
    get 'connect' => 'instagram#connect', as: 'ig_connect'
    get 'callback' => 'instagram#callback', as: 'ig_callback'
  end

  scope 'db' do
    get 'connect' => 'dropbox#connect', as: 'db_connect'
    get 'callback' => 'dropbox#callback', as: 'db_callback'
  end

  get 'sync' => 'pages#sync'
end
