
import React from 'react'
import ReactQuill from 'react-quill';

import 'react-quill/dist/quill.snow.css';

const toolbarOptions = [
    [{ 'font': [] }],
    [{ 'header': [1, 2, 3, 4, false] }],
    [{ 'align': [] }],
    ['blockquote', 'code-block'],
    ['bold', 'italic', 'underline', 'strike'],        // toggled buttons

    [{ 'list': 'ordered' }, { 'list': 'bullet' }],
    [{ 'script': 'sub' }, { 'script': 'super' }],      // superscript/subscript
    [{ 'indent': '-1' }, { 'indent': '+1' }],          // outdent/indent
    [{ 'direction': 'rtl' }],                         // text direction

    [{ 'color': [] }, { 'background': [] }],          // dropdown with defaults from theme
];



function TextEditor(props, ref) {

    return (
        <ReactQuill modules={{
            toolbar: toolbarOptions
        }} theme="snow" {...props} ref={ref} />
    )
}

export default React.forwardRef(TextEditor)