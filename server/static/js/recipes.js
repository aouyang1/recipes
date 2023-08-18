store = {
    listSubItems: null, // list of sub items to render and update
}

function clearListItems() {
    d3.select("#list-items").selectAll("*").remove();
}

function clearListSubItems() {
    d3.select("#list-sub-items").selectAll("*").remove();
}

function clearRecipeUpdate() {
    d3.select("#title-sub-item").selectAll("*").remove();
    d3.select("#badges-recipe-tags").selectAll("*").remove();
    d3.select("#table-recipe-ingredients").selectAll("*").remove();
    d3.select("#button-save-recipe").selectAll("*").remove();
}

d3.select("#sidebar-manage-events")
    .on("click", function() {
        d3.select("#list-items-title").text("Events");
        getRecipeEvents();
    })

d3.select("#sidebar-manage-tags")
    .on("click", function() {
        d3.select("#list-items-title").text("Tags");
        getTags();
    })

d3.select("#sidebar-manage-ingredients")
    .on("click", function() {
        d3.select("#list-items-title").text("Ingredients");
        getIngredients();
    })


d3.select("#list-items-title").text("Events");
getRecipeEvents();