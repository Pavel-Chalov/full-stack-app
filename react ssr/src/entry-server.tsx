import ReactDOMServer from 'react-dom/server'
import { App } from './app'
import { StaticRouter } from "react-router-dom/server";

export function render(url: string) {
  const html = ReactDOMServer.renderToString(
    <StaticRouter location={url}><App /></StaticRouter>
  )
  return { html }
}
