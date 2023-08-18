function getIngredients() {
    d3.request("/ingredients")
        .get(function(error, data) {
            if (error) throw error;
            ingredients = JSON.parse(data.response);
            if (ingredients) {
                list = d3.select("#list-items");
                clearListItems()
                clearRecipeUpdate();
                clearListSubItems();

                /*
                <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
                    <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
                </a>
                */
                list.selectAll("a")
                    .data(ingredients)
                    .enter()
                    .append("a")
                        .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
                        .attr("id", d => d.id)
                        .attr("data-bs-toggle", "list")
                        .on("click", (_, d) => {
                            clearRecipeUpdate();
                            clearListSubItems();
                        })
                    .append("d")
                        .attr("class", "col-10 mb-1 small")
                        .text(d => d.name);
            }
        })
}

