import axios from "axios";

async function fetchData(start: number, end: number) {
  const promises: any = [];

  for (let i = start; i <= end; i++) {
    promises.push(axios.get(`https://pokeapi.co/api/v2/pokemon/${i}/`));

    if (promises.length >= 10) {
      // Wait for all 10 requests to complete before making more
      let res = await Promise.all(promises);
      promises.length = 0; // Clear the array for the next batch
      res.forEach((el) => {
        console.log(el.data.base_experience);
      });
    }
  }

  // Wait for any remaining requests to complete
  await Promise.all(promises);
}

fetchData(1, 100);
