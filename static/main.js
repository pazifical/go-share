function toParent() {
    currentPath = currentPath.substring(0, currentPath.lastIndexOf("/"));
    updateDirEntries(currentPath, true);
}

async function addToPath(name) {
    currentPath = `${currentPath}/${name}`;
    await updateDirEntries(currentPath, true);
}

async function updateDirEntries(currentPath, changePath) {
    const response = await fetch(`/api/list/${currentPath}`);
    const text = await response.text();
    document.getElementById("entries").innerHTML = await text;

    if (changePath === true) {
        window.location.search = `?path=${currentPath}`;
    }
}


const params = new URLSearchParams(window.location.search);
currentPath = params.get("path");
if (!currentPath) {
    currentPath = ".";
}

updateDirEntries(currentPath, false);