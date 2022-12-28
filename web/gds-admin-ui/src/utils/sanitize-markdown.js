import DOMPurify from 'dompurify';


export default function sanitizeMarkdown(markdown){
    return DOMPurify.sanitize(markdown)
}