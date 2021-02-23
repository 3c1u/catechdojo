with
    rarity as (select rand() as rarity),
    characters_with_rarity as (
        select
            characters.id as id,
            characters.name as name,
            rate
        from characters
        join rarities
        on characters.rarity_id = rarities.id
    )
select id, name, rate, rarity from
    characters_with_rarity, rarity
where rate < rarity
order by rate desc, rand();
