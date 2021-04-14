import dotenv from "dotenv";
import main from "./main";

dotenv.config();

(async () => {
    await main();
})();
