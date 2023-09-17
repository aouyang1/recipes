function getIngredients() {
    d3.request("/ingredients")
        .get(function(error, data) {
            if (error) throw error;
            ingredients = JSON.parse(data.response);
            store.listItems = ingredients;
            if (ingredients) {
                list = d3.select("#list-items");
                clearListItems()
                clearRecipeUpdate();
                clearListSubItems();

                renderListNewButton(createIngredient);
                renderListItems();

                d3.select("#input-search-query")
                    .on("keyup", function() {
                        clearListItems();
                        renderListNewButton(createIngredient);
                        renderListItems();
                    });
            }
        })
}

function createIngredient(ingredient_name) {
    req = {"name": ingredient_name};

    d3.request("/ingredient")
        .post(JSON.stringify(req), function(error, data) {
            if (error) {
                console.log(error);
                return
            }
            if (data.response) {
                tag = JSON.parse(data.response); 
                store.listItems.push(tag);
                renderListItems();
            }
        });
}