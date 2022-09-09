select *
from metrics
where idrec in (select max(idrec)  
from metrics
group by id, mtype)