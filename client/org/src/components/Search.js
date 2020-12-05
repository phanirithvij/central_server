import { AlgoliaProvider } from "leaflet-geosearch";
import React, { useEffect, useState } from "react";
import styles from "./Search.module.css";

const provider = new AlgoliaProvider();
function Search(props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);

  const [run, setRun] = useState(true);
  useEffect(() => {
    if (run) {
      if (query === "") return;
      //   console.log("don't run for 300 ms");
      setRun(false);
      setTimeout(() => {
        // console.log("run enabled");
        setRun(true);
      }, 300);
      //   console.log("Search", query);
      provider
        .search({ query })
        .then((results) => setResults(results.slice(0, 10)))
        .catch((err) => console.error(err));
      return;
    }
    // exhaustive dependencies not needed
    // eslint-disable-next-line
  }, [query]);

  const handleClick = (x) => {
    if (props.selectCallback) props.selectCallback(x);
    setQuery(x.label);
    setTimeout(() => {
      setResults([]);
    }, 200);
  };

  const search = (nquery) => {
    //   no duplicate searchs or empty searchs
    if (nquery === "" || nquery === query) return;
    setQuery(nquery);
  };

  return (
    <div className={styles.search}>
      <form>
        <input
          type="text"
          placeholder="Search address"
          value={query}
          onChange={(e) => search(e.target.value)}
        />
      </form>

      <div className={styles.result}>
        {results.map((result, idx) => (
          <div key={idx} onClick={() => handleClick(result)}>
            {result.label}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Search;
