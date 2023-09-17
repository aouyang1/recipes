function getTags() {
    d3.request("/tags")
        .get(function(error, data) {
            if (error) throw error;
            tags = JSON.parse(data.response);
            store.listItems = tags;
            if (tags) {
                clearListItems()
                clearRecipeUpdate();
                clearListSubItems();

                renderListNewButton(createTag);
                renderListItems();

                d3.select("#input-search-query")
                    .on("keyup", function() {
                        clearListItems();
                        renderListNewButton(createTag);
                        renderListItems();
                    });

            }
        })
}

function createTag(tag_name) {
    req = {"name": tag_name};

    d3.request("/tag")
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