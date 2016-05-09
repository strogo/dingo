package model

import (
	"github.com/dinever/dingo/app/model/sql_builder"
)

const schema = `
CREATE TABLE IF NOT EXISTS
posts (
  id                 integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  title              varchar(150) NOT NULL,
  slug               varchar(150) NOT NULL,
  markdown           text,
  html               text,
  image              text,
  featured           BOOLEAN,
  page               BOOLEAN,
  allow_comment      BOOLEAN,
  published          BOOLEAN,
  comment_num        integer NOT NULL DEFAULT '0',
  language           varchar(6) NOT NULL DEFAULT 'en_US',
  meta_title         varchar(150),
  meta_description   varchar(200),
  created_at         datetime NOT NULL,
  created_by         integer NOT NULL,
  updated_at         datetime,
  updated_by         integer,
  published_at       datetime,
  published_by       integer
);

CREATE TABLE IF NOT EXISTS
tokens (
  value       varchar(40) NOT NULL PRIMARY KEY,
  user_id     integer UNIQUE,
  created_at  datetime,
  expired_at  datetime
);

CREATE TABLE IF NOT EXISTS
users (
  id               integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name             varchar(150) NOT NULL,
  slug             varchar(150) NOT NULL,
  password         varchar(60) NOT NULL,
  email            varchar(254) NOT NULL,
  image            text,
  cover            text,
  bio              varchar(200),
  website          text,
  location         text,
  accessibility    text,
  status           varchar(150) NOT NULL DEFAULT 'active',
  language         varchar(6) NOT NULL DEFAULT 'en_US',
  meta_title       varchar(150),
  meta_description varchar(200),
  last_login       datetime,
  created_at       datetime NOT NULL,
  created_by       integer NOT NULL,
  updated_at       datetime,
  updated_by       integer
);

CREATE TABLE IF NOT EXISTS
categories (
  id                integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name              varchar(150) NOT NULL,
  slug              varchar(150) NOT NULL,
  description       varchar(200),
  parent_id         integer,
  meta_title        varchar(150),
  meta_description  varchar(200),
  created_at        datetime NOT NULL,
  created_by        integer NOT NULL,
  updated_at        datetime,
  updated_by        integer
);

CREATE TABLE IF NOT EXISTS
tags (
  id                integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name              varchar(150) NOT NULL,
  slug              varchar(150) NOT NULL,
  description       varchar(200),
  image             text,
  hidden            boolean NOT NULL DEFAULT 0,
  parent_id         integer,
  meta_title        varchar(150),
  meta_description  varchar(200),
  created_at        datetime NOT NULL,
  created_by        integer NOT NULL,
  updated_at        datetime,
  updated_by        integer
);

CREATE TABLE IF NOT EXISTS
comments (
  id            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  post_id       varchar(150) NOT NULL,
  author        varchar(150) NOT NULL,
  author_email  varchar(150) NOT NULL,
  author_url    varchar(200) NOT NULL,
  author_ip     varchar(100) NOT NULL,
  created_at    datetime NOT NULL,
  content       text NOT NULL,
  approved      tinyint NOT NULL DEFAULT '0',
  agent         varchar(255) NOT NULL,
  type          varchar(20),
  parent        integer,
  user_id       integer
);

CREATE TABLE IF NOT EXISTS
posts_tags (
  id       integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  post_id  integer NOT NULL,
  tag_id   integer NOT NULL
);

CREATE TABLE IF NOT EXISTS
posts_categories (
  id           integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  post_id      integer NOT NULL,
  category_id  integer NOT NULL
);

CREATE TABLE IF NOT EXISTS
settings (
  id          integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  key         varchar(150) NOT NULL,
  value       varchar(20) NOT NULL,
  type        varchar(150) NOT NULL DEFAULT 'core',
  created_at  datetime NOT NULL,
  created_by  integer NOT NULL,
  updated_at  datetime,
  updated_by  integer
);

CREATE TABLE IF NOT EXISTS
roles (
  id           integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name         varchar(150) NOT NULL,
  description  varchar(200),
  created_at   datetime NOT NULL,
  created_by   integer NOT NULL,
  updated_at   datetime,
  updated_by   integer
);

CREATE TABLE IF NOT EXISTS
messages (
  id           integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  type         varchar(20) NOT NULL,
  data         text NOT NULL,
  is_read      boolean NOT NULL default 0,
  created_at   datetime NOT NULL
);
`

// Posts
var postCountSelector = SQL.Select(`count(*)`).From(`posts`)

var stmtGetPostsCountByTag = postCountSelector.Copy().From(`posts, posts_tags`).Where(`posts_tags.post_id = posts.id`, `posts_tags.tag_id = ?`, `posts.published`).SQL()

const stmtDeletePostById = `DELETE FROM posts WHERE id = ?`

//PostTags
const stmtDeletePostTagsByPostId = `DELETE FROM posts_tags WHERE post_id = ?`
const stmtInsertPostTag = `INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)`

// Comments
var commentCountSelector = SQL.Select(`count(*)`).From(`comments`)
var stmtGetAllCommentCount = commentCountSelector.SQL()
var commentSelector = SQL.Select(`id, post_id, author, author_email, author_url, created_at, content, approved, agent, parent, user_id`).From(`comments`)
var stmtGetAllCommentList = commentSelector.Copy().OrderBy(`created_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetApprovedCommentList = commentSelector.Copy().Where(`approved = 1`).OrderBy(`created_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetCommentById = commentSelector.Copy().Where(`id = ?`).SQL()
var stmtGetApprovedCommentListByPostId = commentSelector.Copy().Where(`post_id = ?`, `approved = 1`).OrderBy(`created_at DESC`).SQL()

const stmtInsertComment = `INSERT OR REPLACE INTO comments (id, post_id, author, author_email, author_url, author_ip, created_at, content, approved, agent, parent, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const stmtDeleteCommentById = `DELETE FROM comments WHERE id = ?`

// Users
const stmtGetUserById = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE id = ?`
const stmtGetUserBySlug = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE slug = ?`
const stmtGetUserByName = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE name = ?`
const stmtGetUserByEmail = `SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE email = ?`
const stmtGetHashedPasswordByEmail = `SELECT password FROM users WHERE email = ?`
const stmtGetUsersCountByEmail = `SELECT count(*) FROM users where email = ?`
const stmtInsertUser = `INSERT INTO users (id, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const stmtUpdateUser = `UPDATE users SET name = ?, slug = ?, email = ?, image = ?, cover = ?, bio = ?, website = ?, location = ?, updated_at = ?, updated_by = ? WHERE id = ?`

// RoleUser
const stmtInsertRoleUser = `INSERT INTO roles_users (id, role_id, user_id) VALUES (?, ?, ?)`

// Tokens
const stmtGetTokenByValue = `SELECT value, user_id, created_at, expired_at FROM tokens WHERE value = ?`
const stmtUpdateToken = `INSERT OR REPLACE INTO tokens (value, user_id, created_at, expired_at) VALUES (?, ?, ?, ?)`

// Tags
const stmtGetAllTags = `SELECT id, name, slug FROM tags`
const stmtGetTags = `SELECT tag_id FROM posts_tags WHERE post_id = ?`
const stmtGetTagById = `SELECT id, name, slug FROM tags WHERE id = ?`
const stmtGetTagBySlug = `SELECT id, name, slug, hidden FROM tags WHERE slug = ?`
const stmtInsertTag = `INSERT INTO tags (id, name, slug, created_at, created_by, updated_at, updated_by, hidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
const stmtUpdateTag = `UPDATE tags SET name = ?, slug =?, updated_at = ?, updated_by = ?, hidden = ? WHERE id = ?`
const stmtDeleteOldTags = `DELETE FROM tags WHERE id IN (SELECT id FROM tags EXCEPT SELECT tag_id FROM posts_tags)`

// Settings
const stmtGetSettingByKey = `SELECT id, key, value, type, created_at, created_by from settings where key = ?`
const stmtGetSettingsByType = `SELECT id, key, value, type, created_at, created_by from settings where type = ?`
const stmtUpdateSetting = `INSERT OR REPLACE INTO settings (id, key, value, type, created_at, created_by) VALUES ((SELECT id FROM settings WHERE key = ?), ?, ?, ?, ?, ?)`

// Messages
var messageSelector = SQL.Select(`id, type, data, is_read, created_at`).From(`messages`)

//var stmtGetUnreadMessages = messageSelector.Copy().Where(`is_read = 0`).OrderBy(`created_at DESC`).Limit(`?`).Offset(`?`).SQL()
var stmtGetUnreadMessages = messageSelector.Copy().Where(`is_read = 0`).OrderBy(`created_at DESC`).Limit(`10`).Offset(`0`).SQL()
