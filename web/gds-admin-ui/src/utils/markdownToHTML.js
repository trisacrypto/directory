import {marked} from 'marked';

export default function markdownToHTML(markdown){
    return marked.parse(markdown)
}