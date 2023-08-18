function getRecipeEvents() {
    d3.request("/recipe_events")
        .get(function(error, data) {
            if (error) throw error;
            recipe_events = JSON.parse(data.response);
            list = d3.select("#list-items");
            list.selectAll("*").remove();

            /*
            <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
                <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
            </a>
            */
            list.selectAll("a")
                .data(recipe_events)
                .enter()
                .append("a")
                    .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
                    .attr("id", d => d.id)
                    .attr("data-bs-toggle", "list")
                    .on("click", (_, d) => {
                        getRecipesByRecipeEventID(d.id);
                    })
                .append("d")
                    .attr("class", "col-10 mb-1 small")
                    .text(d => d.title);
        })
}


function getRecipesByRecipeEventID(recipe_event_id) {
    req = {"recipe_event_id": recipe_event_id};
    params = new URLSearchParams(req).toString();
    d3.request("/recipes?" + params)
        .get(function(error, data) {
            if (error) throw error;
            recipes = JSON.parse(data.response);
            console.log(recipes);

            list = d3.select("#list-sub-items");
            list.selectAll("*").remove();

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
                    .attr("id", "dropdownFormName")
                    .attr("placeholder", "recipe name");

            btnForm
                .append("div")
                    .attr("class", "form-group py-1")
                .append("input")
                    .attr("type", "text")
                    .attr("class", "form-control")
                    .attr("id", "dropdownFormLink")
                    .attr("placeholder", "https://recipelink.com");

            btnForm
                .append("button")
                    .attr("type", "submit")
                    .attr("class", "btn btn-primary py-2")
                    .text("Create")
                .on("click", (_, d) => {
                    console.log(d);
                    // if successful creation of recipe to event append to recipes
                    // variable and focus on the newly created one
                    recipe_name = d3.select("#dropdownFormName").property('value');
                    recipe_variant = d3.select("#dropdownFormLink").property('value'); 
                    CreateLinkEventToRecipe(d, recipe_name, recipe_variant)
                });

            AppendSubList(recipes);
        })
}

function CreateLinkEventToRecipe(recipe_event_id, recipe_name, recipe_variant) {
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
            recipe = JSON.parse(data.response); 
            console.log(recipe);
            AppendSubList([recipe])
        });
}

function AppendSubList(items) {
    /*
    <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
        <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
    </a>
    */

    list = d3.select("#list-sub-items");
    list.selectAll("a")
        .data(items)
        .enter()
        .append("a")
            .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
            .attr("id", d => d.name + ":" + d.variant)
            .attr("data-bs-toggle", "list")
            .on("click", (_, recipe) => {
                console.log(recipe);
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
                        console.log(recipe);
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
                        console.log(recipe);
                    });

                badges = d3.select("#badges-recipe-tags");
                /*
                <div class="input-group mb-3">
                  <div class="input-group-prepend">
                    <span class="input-group-text" id="basic-addon1">@</span>
                  </div>
                  <input type="text" class="form-control" placeholder="Username" aria-label="Username" aria-describedby="basic-addon1">
                </div>
                */
                tagInput = badges.append("div")
                    .attr("class", "col-6 input-group input-group-sm p-2");

                tagInput.append("div")
                        .attr("class", "input-group-prepend")
                    .append("span")
                        .attr("class", "input-group-text")
                        .text("Tags");
                tagInput.append("input")
                    .attr("class", "form-control")
                    .attr("type", "text");

                /*
                <span class="badge badge-pill badge-info">Info</span>
                */
                badges.selectAll("a")
                    .data(recipe.tags)
                    .enter()
                    .append("a")
                        .attr("class", "p-1")
                        .attr("id", d => "badge-" + d)
                    .append("span")
                        .attr("class", "badge rounded-pill bg-info text-dark")
                        .text(d => d)
                    .append("i")
                        .attr("class", "icon-remove px-1")
                        .attr("style", "color:red")
                        .on("click", (i, d) => {
                            d3.select("#badge-"+d).remove();
                            idx = recipe.tags.indexOf(d);
                            recipe.tags.splice(idx, 1);
                            console.log(recipe);
                        });


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
                        console.log("saving");
                        d3.request("/recipe")
                            .send("PUT", JSON.stringify(recipe), function(error, data) {
                                if (error) {
                                    console.log(error);
                                    return
                                }
                                console.log(recipe);
                            });
                    });

            })
        .append("d")
            .attr("class", "col-10 mb-1 small")
            .text(d => d.name);
}

getRecipeEvents();