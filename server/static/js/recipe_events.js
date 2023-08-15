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
            list = d3.select("#list-sub-items");
            list.selectAll("*").remove();

            /*
            <div class="btn-group">
              <button class="btn btn-secondary btn-sm dropdown-toggle" type="button" data-toggle="dropdown">
                Create
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
                    .attr("class", "btn-group");

            btnGroup
                .append("button")
                    .attr("class", "btn btn-secondary btn-sm dropdown-toggle")
                    .attr("type", "button")
                    .attr("data-bs-toggle", "dropdown")
                    .text("New Recipe");
 
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
                    .attr("class", "btn btn-primary")
                    .text("Create")
                .on("click", (_, d) => {
                    console.log(d);
                    // if successful creation of recipe to event append to recipes
                    // variable and focus on the newly created one
                });

            /*
            <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
                <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
            </a>
            */
            list.selectAll("a")
                .data(recipes)
                .enter()
                .append("a")
                    .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
                    .attr("id", d => d.name + ":" + d.variant)
                    .attr("data-bs-toggle", "list")
                    .on("click", (_, d) => {
                        console.log(d);
                    })
                .append("d")
                    .attr("class", "col-10 mb-1 small")
                    .text(d => d.name);
        })
}

getRecipeEvents();