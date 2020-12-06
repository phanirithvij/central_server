import { Puff, useLoading } from "@agney/react-loading";
import { OpenStreetMapProvider } from "leaflet-geosearch";
import React, { useState } from "react";
import styles from "./Search.module.css";

// eslint-disable-next-line no-extend-native
String.prototype.trimLeft =
  String.prototype.trimLeft ||
  function () {
    var start = -1;
    while (this.charCodeAt(++start) < 33);
    return this.slice(start, this.length);
  };

export async function Address(latlong) {
  return await provider.search({ query: latlong.toString() });
}

const provider = new OpenStreetMapProvider();
function Search(props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);

  const [run, setRun] = useState(true);

  const handleClick = (x) => {
    if (props.selectCallback) props.selectCallback(x);
    setQuery(x.label);
    setTimeout(() => {
      setSearching(undefined);
      setResults([]);
    }, 200);
  };

  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  const [searching, setSearching] = useState();

  const search = (nquery) => {
    // console.log(nquery + "x");
    // no empty searchs and < 3 searches
    if (nquery === "" || nquery.length < 3) {
      if (results.length !== 0) {
        setSearching(undefined);
        setResults([]);
      }
      setQuery(nquery.trimLeft());
      return;
    }
    // no duplicate searchs
    // if (nquery.trim() === query.trim()) {
    //   setQuery(nquery.trimLeft());
    //   return;
    // }
    setQuery(nquery.trimLeft());
    if (run) {
      //   console.log("don't run for 300 ms");
      setRun(false);
      setTimeout(() => {
        // console.log("run enabled");
        setRun(true);
      }, 300);
      // console.log("Search", nquery.trimLeft());
      setSearching(true);
      provider
        .search({ query: nquery.trimLeft() })
        .then((results) => {
          setSearching(false);
          setResults(results.slice(0, 10));
        })
        .catch((err) => console.error(err));
      return;
    }
    // exhaustive dependencies not needed
    // eslint-disable-next-line
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
        {searching !== undefined &&
          (searching ? (
            <section {...containerProps}>{indicatorEl}</section>
          ) : results.length === 0 ? (
            "No results found"
          ) : (
            ""
          ))}
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
