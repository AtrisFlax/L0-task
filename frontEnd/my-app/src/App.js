import React, { useRef, useState } from "react";
import './App.css';

function App() {
  const baseURL = "http://localhost:8000";

  const get_uuid = useRef(null);

  const post_title = useRef(null);

  const [getResult, setGetResult] = useState(null);
  const [postResult, setPostResult] = useState(null);

  const fortmatResponse = (res) => {
    return JSON.stringify(res, null, 2);
  }

  async function getDataByUUID() {
    const uuid = get_uuid.current.value;

    if (uuid) {
      try {
        const res = await fetch(`${baseURL}/task/${uuid}`);

        if (!res.ok) {
          const message = `An error has occured: ${res.status} - ${res.statusText}`;
          throw new Error(message);
        }

        const data = await res.json();

        const result = {
          data: data,
          status: res.status,
          statusText: res.statusText,
          headers: {
            "Content-Type": res.headers.get("Content-Type"),
            "Content-Length": res.headers.get("Content-Length"),
          },
        };

        setGetResult(fortmatResponse(result));
      } catch (err) {
        setGetResult(err.message);
      }
    }
  }

  async function postData() {
    const postData = post_title.current.value;

    try {

      console.log(postData)

      const res = await fetch(`${baseURL}/task`, {
        credentials: "include",
        method: "post",
        headers: {
          "Content-Type": "application/json",
        },
        body: postData,
      });
          
      if (!res.ok) {
        const message = `An error has occured: ${res.status} - ${res.statusText}`;
        throw new Error(message);
      }

      const data = await res.json();

      const result = {
        status: res.status + "-" + res.statusText,
        headers: {
          "Content-Type": res.headers.get("Content-Type"),
          "Content-Length": res.headers.get("Content-Length"),
        },
        data: data,
      };

      setPostResult(fortmatResponse(result));
    } catch (err) {
      setPostResult(err.message);
    }
  }

  const clearGetOutput = () => {
    setGetResult(null);
  }

  const clearPostOutput = () => {
    setPostResult(null);
  }

  return (
    <div id="app" className="container my-3">
      <h3>React Fetch example</h3>

      <div className="card mt-3">
        <div className="card-header">GET</div>
        <div className="card-body">
          <div className="input-group input-group-sm">
           
            <input type="text" ref={get_uuid} className="form-control ml-2" placeholder="uuid" />
            <div className="input-group-append">
              <button className="btn btn-sm btn-primary" onClick={getDataByUUID}>Get by Id</button>
            </div>

            <button className="btn btn-sm btn-warning ml-2" onClick={clearGetOutput}>Clear</button>
          </div>   
          
          { getResult && <div className="alert alert-secondary mt-2" role="alert"><pre>{getResult}</pre></div> }
        </div>
      </div>

      <div className="card mt-3">
        <div className="card-header">POST</div>
        <div className="card-body">
          <button className="btn btn-sm btn-primary" onClick={postData}>Post Data</button>
          <button className="btn btn-sm btn-warning ml-2" onClick={clearPostOutput}>Clear</button>
          { postResult && <div className="alert alert-secondary mt-2" role="alert"><pre>{postResult}</pre></div> }
          <div className="form-group">
            <textarea type="text" className="form-control" rows="62" ref={post_title} placeholder="Json Data" />
          </div>

          
        </div>
      </div> 
    </div>
  );
}

export default App;
