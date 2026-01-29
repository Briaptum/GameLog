-- Remove testimonials for Guthrie Zaring and Alissa Meriwether
DELETE FROM testimonials WHERE for_user IN ('Guthrie Zaring', 'Alissa Meriwether');
