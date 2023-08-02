DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id` INT AUTO_INCREMENT NOT NULL,
  `title` TEXT NULL,
  `content` TEXT NULL,
  `created_at` TIMESTAMP  DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO `posts`
  (title, content)
VALUES
  ('テストタイトル１', '投稿１です。投稿１です。投稿１です。投稿１です。投稿１です。投稿１です。投稿１です。'),
  ('テストタイトル2', '投稿2です。投稿2です。投稿2です。投稿2です。投稿2です。投稿2です。投稿2です。'),
  ('テストタイトル3', '投稿3です。投稿3です。投稿3です。投稿3です。投稿3です。投稿3です。投稿3です.'),
  ('テストタイトル4', '投稿4です。投稿4です。投稿4です。投稿4です。投稿4です。投稿4です。投稿4です。'),
  ('テストタイトル5', '投稿５です。投稿５です。投稿５です。投稿５です。投稿５です。投稿５です。投稿５です。'),
  ('テストタイトル6', '投稿6です。投稿6です。投稿6です。投稿6です。投稿6です。投稿6です。投稿6です。');
