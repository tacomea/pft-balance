NAME = pft-balance

FOOD_PATH	= food
MENU_PATH	= menu
STORE_PATH	= store
USER_PATH	= user

all: $(NAME)

$(NAME):
	@make -C $(FOOD_PATH)
	@make -C $(MENU_PATH)
	@make -C $(STORE_PATH)
	@make -C $(USER_PATH)

blog_server:
	go run blog/blog_server/server.go

blog_client:
	go run blog/blog_client/client.go