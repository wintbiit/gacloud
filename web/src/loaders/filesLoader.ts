import {defer} from "react-router-dom";
import {listFiles} from "../api/files.ts";
const filesLoader = async ({params}: any) => {
    const path = params["*"];

    console.log("Loading files from", path);

    const files = await listFiles(path).then(files => {
        // sort: has / suffix first
        return files.sort((a, b) => {
            if (a.path.endsWith("/") && !b.path.endsWith("/")) {
                return -1;
            } else if (!a.path.endsWith("/") && b.path.endsWith("/")) {
                return 1;
            } else {
                return a.path.localeCompare(b.path);
            }
        });
    });

    return defer({files: files, path: path})
}

export default filesLoader;