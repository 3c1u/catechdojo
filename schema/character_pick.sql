with characters_with_sum as (select name,
            rate,
        sum(rate) over(order by id) as sum_rate,
        sum(rate) over() as total_rate
    from characters)
select * from
    characters_with_sum
     where
    sum_rate > total_rate * rand()
     order by rate desc, rand()
     limit 1;
