# omo-msa-third
Micro Service Agent - third
合作伙伴或者渠道

MICRO_REGISTRY=consul micro call omo.msa.third FavoriteService.AddOne '{"name":"John", "owner":"11111", "remark":"test1", "type":2, "cover":"hhhhhh"}'
MICRO_REGISTRY=consul micro call omo.msa.third FavoriteService.RemoveOne '{"uid":"5f0fffef0d57c9d90026b782", "owner":"11111"}'
