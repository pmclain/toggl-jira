import dotenv from "dotenv";
import main from "./main.js";

dotenv.config();

(async () => {
    await main();
})();
