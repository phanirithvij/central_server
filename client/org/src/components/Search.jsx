import { Puff, useLoading } from "@agney/react-loading";
import { OpenStreetMapProvider } from "leaflet-geosearch";
import React, { useState } from "react";
import styles from "./Search.module.css";

// eslint-disable-next-line no-extend-native
String.prototype.trimLeft =
  String.prototype.trimLeft ||
  // https://stackoverflow.com/a/1593909/8608146
  function () {
    var start = -1;
    while (this.charCodeAt(++start) < 33);
    return this.slice(start, this.length);
  };

/**
 * Uses the provider to query the address from the latlong
 * @param {[number, number]} latlong
 */
export async function Address(latlong) {
  return await provider.search({ query: latlong.toString() });
}

// Using the openstreetmap provider
// Dev notes: Check https://smeijer.github.io/leaflet-geosearch/providers/algolia
// for more providers
const provider = new OpenStreetMapProvider();

function Search(props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);

  const [run, setRun] = useState(true);

  const handleClick = (x) => {
    if (props.selectCallback) props.selectCallback(x);
    setQuery(x.label);
    setTimeout(() => {
      // reset search results
      setResults([]);
      // don't show search results
      setSearching(undefined);
    }, 200);
  };

  // https://reactjsexample.com/simple-and-accessible-loading-indicators-with-react/
  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  // initially undefined
  const [searching, setSearching] = useState();

  const search = (nquery) => {
    // no empty searchs and < 3 searches
    if (nquery === "" || nquery.length < 3) {
      if (results.length !== 0) {
        // show nothing
        setSearching(undefined);
        setResults([]);
      }
      setQuery(nquery.trimLeft());
      return;
    }
    setQuery(nquery.trimLeft());
    if (run) {
      setRun(false);
      setTimeout(() => {
        setRun(true);
      }, 300);
      setSearching(true);
      provider
        .search({ query: nquery.trimLeft() })
        .then((results) => {
          setSearching(false);
          setResults(results.slice(0, 10));
        })
        .catch((err) => {
          setSearching(false);
          // this sometimes fails for no reason
          // issue inside the leaflet-geosearch
          console.error(err);
        });
      return;
    }
  };

  return (
    <div className={styles.search}>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          search(query);
        }}
      >
        <input
          type="text"
          placeholder="Search address"
          value={query}
          onChange={(e) => search(e.target.value)}
        />
      </form>

      {searching !== undefined && (
        <div className={styles.result}>
          {searching ? (
            <section {...containerProps}>{indicatorEl}</section>
          ) : results.length === 0 ? (
            "No results found"
          ) : (
            // if results are available
            results.map((result, idx) => (
              <div key={idx} onClick={() => handleClick(result)}>
                {result.label}
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}

export default Search;
