select *
from metrics
where 
id = 'GetSet312'
order by idrec DESC

-- select id, mtype, delta, value
-- 						from metrics
-- 						where idrec in (select max(idrec)  
-- 						from metrics
-- 						group by id)

-- TRUNCATE TABLE metrics

-- select count(*)
-- from metrics