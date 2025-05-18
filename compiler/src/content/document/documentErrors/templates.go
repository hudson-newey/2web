package documentErrors

import "time"

func errorHtmlSource() string {
	currentTime := time.Now().Format(time.DateTime)

	return `
        <div class="__2_error_overlay">
            <div class="__2_error_container">
            	<div class="__2_error_header">
                    <h1 class="__2_error_title">
                        <strong>Compiler Error</strong>
                    </h1>
                    <time class="__2_error_time">` + currentTime + `</time>
                </div>

                <div class="__2_error_overview">
                    <ul>
                        <li>{{ len . }} Errors</li>
                        <li>0 Warnings</li>
                    </ul>
                </div>

                <div class="__2_error_list">
                    {{range .}}
                        <div class="__2_error">
                            <h2 class="__2_error_file">{{.FilePath}}</h2>
                            <pre class="__2_error_message">{{.Message}}</pre>
                        </div>
                    {{end}}
                </div>
            </div>

            <style>
                .__2_error_overlay {
                    position: fixed;
                    inset: 0 0;
                    background-color: rgba(100, 100, 100, 0.5);
                }

                .__2_error_container {
                    position: fixed;
                    inset: 5%;
                    padding: 2rem 5rem;
                    z-index: 5000;

                    background: radial-gradient(circle, rgba(5, 0, 0, 0.9) 0%, rgba(20, 0, 0, 0.98) 99%, rgba(40, 0, 0, 1) 100%);
                    background-color: rgba(20, 0, 0);
                    border-radius: 1rem;

                    box-shadow: 0 0.5rem 10rem rgba(0, 0, 0);

                    font-family: sans-serif;

                    overflow-y: auto;
                }

                .__2_error_header {
                    border-bottom: 1px solid #fff;
                    padding-bottom: 1rem;
                }

                .__2_error_title {
                    font-weight: bold;
                    font-size: 2.2rem;
                    margin-bottom: 0.5rem;

                    color: #fff;
                }

                .__2_error_time {
                    font-size: 1rem;
                    color: #fff;
                    font-style: italic;
                }

                .__2_error_overview {
                    display: flex;
                    margin-bottom: 2rem;
                    border-bottom-left-radius: 0.5rem;
                    border-bottom-right-radius: 0.5rem;

                    background-color: rgba(140, 30, 40, 0.7);
                    color: white;

                    font-size: 1em;
                    letter-spacing: 1px;

                    & > ul {
                        padding-left: 1.5rem;
                        list-style-type: none;

                        li {
                            display: inline-block;
                            margin-right: 2rem;
                        }
                    }
                }

                .__2_error_list {
                    & > .__2_error {
                        color: white;
                        margin-bottom: 1.5rem;
                        padding: 1rem;
                        background-color: rgba(80, 20, 20, 0.1);
                        border-radius: 0.5rem;
                        box-shadow: 0 0 0.25rem rgba(149, 38, 43, 0.9);
                    }

                    .__2_error_file {
                        color: #f00;
                        font-weight: bold;
                        margin-bottom: 1rem;
                        margin-top: 0;
                    }

                    .__2_error_message {
                        color: rgb(255, 230, 230);
                        margin: 0;
                        font-size: 1.1rem;
                        font-weight: 400;
                        font-family: monospace;
                        white-space: pre-wrap;
                        line-height: 2;
                        
                        overflow-wrap: break-word;
                    }
                }
            </style>
        </div>
    `
}
