drop table users, user_saved_chapter, user_saved_article, user_saved_words, grammar_chpater, grammar_article, grammar_exersices, grammar_content CASCADE;

DROP INDEX IF EXISTS idx_grammar_articles_chapter_id;
DROP INDEX IF EXISTS idx_grammar_contents_article_id;
DROP INDEX IF EXISTS idx_grammar_contents_level_id;
DROP INDEX IF EXISTS idx_pronunciation_tables_content_id;
DROP INDEX IF EXISTS idx_grammar_comments_content_id;
DROP INDEX IF EXISTS idx_grammar_comments_user_id;
DROP INDEX IF EXISTS idx_grammar_comments_rating;
DROP INDEX IF EXISTS idx_comment_likes_comment;
DROP INDEX IF EXISTS idx_comment_likes_user;
DROP INDEX IF EXISTS idx_grammar_comments_likes;