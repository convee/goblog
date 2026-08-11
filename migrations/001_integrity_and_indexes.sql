-- Apply once to databases created before v1.3.0.
-- Resolve duplicate names/idents/emails before running if the statements fail.
ALTER TABLE `category` ADD UNIQUE KEY `uk_category_name` (`name`);
ALTER TABLE `page` ADD UNIQUE KEY `uk_page_ident` (`ident`);
ALTER TABLE `post` ADD KEY `idx_post_list` (`status`,`is_top`,`updated_at`);
ALTER TABLE `post` ADD KEY `idx_post_category` (`status`,`category_id`,`is_top`,`updated_at`);
ALTER TABLE `post_content` ADD UNIQUE KEY `uk_post_content_post_id` (`post_id`);
ALTER TABLE `tag` ADD UNIQUE KEY `uk_tag_name` (`name`);
ALTER TABLE `user` ADD UNIQUE KEY `uk_user_email` (`email`);
