function getRecipeEvents() {
    d3.request("/recipe_events")
        .get(function(error, data) {
            if (error) throw error;
            recipe_events = JSON.parse(data.response);
            store.listItems = recipe_events;

            clearListItems()
            clearRecipeUpdate();
            clearListSubItems();
            renderRecipeEvents();
        })
}

function renderRecipeEvents() {
    /*
    <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
        <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
    </a>
    */

    list = d3.select("#list-items");
    list.selectAll("a")
        .data(recipe_events)
        .enter()
        .append("a")
            .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
            .attr("id", d => d.id)
            .attr("data-bs-toggle", "list")
            .on("click", (_, d) => {
                clearRecipeUpdate();
                clearListSubItems();
                getRecipesByRecipeEventID(d.id);
            })
        .append("d")
            .attr("class", "col-10 mb-1 small")
            .html((d) => {
                return "<div class=\"row\">"+
                "<div class=\"col-7\">"+d.title+"</div>"+
                "<div class=\"col-1\">"+(d.count > 0 ? d.count : "")+"</div>"+
                "<div class=\"col-4\">"+d.date.slice(0, 10)+"</div>"+
                "</div>"
            });
}

function getRecipesByRecipeEventID(recipe_event_id) {
    req = {"recipe_event_id": recipe_event_id};
    params = new URLSearchParams(req).toString();
    d3.request("/recipes?" + params)
        .get(function(error, data) {
            if (error) throw error;
            recipes = JSON.parse(data.response);
            store.listSubItems = recipes;

            renderNewLinkRecipeButton(recipe_event_id);
            renderSubList();
        })
}

function renderNewLinkRecipeButton(recipe_event_id) {
    /*
    <div class="btn-group">
      <button class="btn btn-secondary btn-sm dropdown-toggle" type="button" data-toggle="dropdown">
        New/Link 
      </button>
      <div class="dropdown-menu">
        <form class="px-4 py-3">
          <div class="form-group">
            <input type="text" class="form-control" id="dropdownFormName" placeholder="recipe name">
          </div>
          <div class="form-group">
            <input type="text" class="form-control" id="dropdownFormLink" placeholder="https://recipelink.com">
          </div>
          <button type="submit" class="btn btn-primary">Create</button>
        </form>
      </div>
    </div>
    */
    list = d3.select("#list-sub-items");

    btnGroup = list.selectAll("div")
        .data([recipe_event_id])
        .enter()
        .append("div")
            .attr("class", "btn-group p-4");

    btnGroup
        .append("button")
            .attr("class", "btn btn-secondary btn-sm dropdown-toggle")
            .attr("type", "button")
            .attr("data-bs-toggle", "dropdown")
            .text("New/Link Recipe");

    btnForm = btnGroup
        .append("div")
            .attr("class", "dropdown-menu")
        .append("form")
            .attr("class", "px-1");

    btnForm
        .append("div")
            .attr("class", "form-group py-1")
        .append("input")
            .attr("type", "text")
            .attr("class", "form-control")
            .attr("id", "dropdownFormRecipeName")
            .attr("placeholder", "recipe name");

    btnForm
        .append("div")
            .attr("class", "form-group py-1")
        .append("input")
            .attr("type", "text")
            .attr("class", "form-control")
            .attr("id", "dropdownFormRecipeLink")
            .attr("placeholder", "https://recipelink.com");

    btnForm
        .append("button")
            .attr("type", "button")
            .attr("class", "btn btn-primary py-2")
            .text("Create")
        .on("click", (_, d) => {
            // if successful creation of recipe to event append to recipes
            // variable and focus on the newly created one
            recipe_name = d3.select("#dropdownFormRecipeName").property('value');
            recipe_variant = d3.select("#dropdownFormRecipeLink").property('value'); 
            createLinkEventToRecipe(d, recipe_name, recipe_variant)
        });
}

function createLinkEventToRecipe(recipe_event_id, recipe_name, recipe_variant) {
    req = {
        "recipe_event_id": recipe_event_id,
        "recipe": {
            "name": recipe_name,
            "variant": recipe_variant,
        },
    };

    d3.request("/recipe")
        .post(JSON.stringify(req), function(error, data) {
            if (error) {
                console.log(error);
                return
            }
            if (data.response) {
                recipe = JSON.parse(data.response); 
                store.listSubItems.push(recipe);
                renderSubList(recipe);
                renderRecipeUpdate(recipe);

                // update internal count of list items of recipe events
                for (i = 0; i < store.listItems.length; i++) {
                    if (store.listItems[i].id == recipe_event_id) {
                        store.listItems[i].count += 1;
                        break;
                    }
                }
                clearListItems();
                renderRecipeEvents();
            }
        });
}

function renderSubList(selectSubItem) {
    /*
    <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
        <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
    </a>
    */

    list = d3.select("#list-sub-items");
    list.selectAll("a").remove();

    console.log(store.listSubItems);
    list.selectAll("a")
        .data(store.listSubItems)
        .enter()
        .append("a")
            .attr("class", d => {
                classStr = "list-group-item list-group-item-action py-3 lh-tight"
                if (selectSubItem != null && selectSubItem.id == d.id) {
                    classStr += " active"
                }
                return classStr;
            })
            .attr("id", d => d.name + ":" + d.variant)
            .attr("data-bs-toggle", "list")
            .on("click", (_, recipe) => {renderRecipeUpdate(recipe)})
        .append("d")
            .attr("class", "col-10 mb-1 small")
            .text(d => d.name);
}

function renderRecipeUpdateTitle(recipe) {
    title = d3.select("#title-sub-item")
        .append("div")
            .attr("class", "border-bottom");

    nameInput = title.append("div")
        .attr("class", "input-group input-group-sm px-2 pt-4 pb-2");

    nameInput.append("div")
            .attr("class", "input-group-prepend")
        .append("span")
            .attr("class", "input-group-text")
            .text("Name");
    nameInput.append("input")
        .attr("class", "form-control")
        .attr("type", "text")
        .property("value", recipe.name)
        .on('change', function() {
            recipe.name = d3.select(this).property('value');
        });

    variantInput = title.append("div")
        .attr("class", "input-group input-group-sm px-2 pb-4");
        
    variantInput.append("div")
            .attr("class", "input-group-prepend")
        .append("span")
            .attr("class", "input-group-text")
            .text("URL");
    variantInput.append("input")
        .attr("class", "form-control")
        .attr("type", "text")
        .property("value", recipe.variant)
        .on('change', function() {
            recipe.variant = d3.select(this).property('value');
        });
}

function renderRecipeTagDropdown(recipe) {
    d3.request("/tags")
        .get(function(error, data) {
            if (error) throw error;
            tags = JSON.parse(data.response);
            if (tags) {
                badges = d3.select("#badges-recipe-tags");

                tagInput = badges.append("div")
                    .attr("class", "col-6 input-group input-group-sm p-2");

                tagInput.append("button")
                    .attr("id", "input-update-tag")
                    .attr("class", "form-control dropdown-toggle")
                    .attr("type", "button")
                    .attr("data-bs-toggle", "dropdown")
                    .attr("aria-haspopup", "true")
                    .attr("aria-expanded", "false")
                    .text("Tags");

                tagInput.append("div")
                    .attr("class", "dropdown-menu")
                    .attr("aria-labelledby", "input-update-tag")
                    .selectAll("a")
                    .data(tags)
                    .enter()
                    .append("a")
                        .attr("class", "dropdown-item")
                        .text(d => d.name)
                        .on("click", function(_, d) {
                            console.log(d);
                            if (recipe.tags.indexOf(d) < 0) {
                                recipe.tags.push(d);
                                renderRecipeUpdateTags(recipe);
                            }
                        });
            }

            renderRecipeUpdateTags(recipe);
            renderRecipeUpdateIngredients(recipe);
        })
}

function renderRecipeUpdateTags(recipe) {
    badges = d3.select("#badges-recipe-tags");
    badges.selectAll("#badge-collection").remove();
    tags = badges.append("div")
        .attr("id", "badge-collection");

    console.log(recipe.tags)
    /*
    <span class="badge badge-pill badge-info">Info</span>
    */
    tags.selectAll("a")
        .data(recipe.tags)
        .enter()
        .append("a")
            .attr("class", "p-1")
            .attr("id", d => "badge-" + d.name)
        .append("span")
            .attr("class", "badge rounded-pill bg-info text-dark")
            .text(d => d.name)
        .append("i")
            .attr("class", "icon-remove px-1")
            .attr("style", "color:red")
            .on("click", (_, d) => {
                console.log(d);
                d3.select("#badge-"+d.name).remove();
                idx = recipe.tags.indexOf(d);
                recipe.tags.splice(idx, 1);
            });
}

function renderRecipeUpdateIngredients(recipe) {
    table = d3.select("#table-recipe-ingredients");
    /*
    <div class="input-group mb-3">
      <div class="input-group-prepend">
        <span class="input-group-text" id="basic-addon1">@</span>
      </div>
      <input type="text" class="form-control" placeholder="Username" aria-label="Username" aria-describedby="basic-addon1">
    </div>
    */
    ingInput = table.append("div")
        .attr("class", "col-6 input-group input-group-sm px-2 pb-2 pt-5");

    ingInput.append("div")
            .attr("class", "input-group-prepend")
        .append("span")
            .attr("class", "input-group-text")
            .text("Ingredients");
    ingInput.append("input")
        .attr("class", "form-control")
        .attr("type", "text");
}

function renderRecipeUpdate(recipe) {
    clearRecipeUpdate();
    renderRecipeUpdateTitle(recipe);
    renderRecipeTagDropdown(recipe);

    /*
    <button type="button" class="btn btn-primary btn-sm">Small button</button>
    */
    btn = d3.select("#button-save-recipe")
        .append("div")
        .attr("class", "col-3 px-2");

    btn.append("button")
        .attr("class", "btn btn-primary btn-sm")
        .text("Save")
        .on("click", (_, d) => {
            d3.request("/recipe")
                .send("PUT", JSON.stringify(recipe), function(error, data) {
                    if (error) {
                        console.log(error);
                        return
                    }
                    console.log(recipe);
                });
        });
}