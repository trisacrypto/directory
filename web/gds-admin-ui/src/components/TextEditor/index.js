import React from 'react'
import MarkdownEditor from "@uiw/react-markdown-editor";
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';

function TextEditor({ value, ...props }, ref) {
    return (
        <>
            <Tabs
                defaultActiveKey="editor"
                transition={false}
                id="editor"
            >
                <Tab eventKey="editor" title="Editor" className='my-1' data-color-mode="light">
                    <MarkdownEditor minHeight='100px'
                        toolbars={['bold', 'italic', 'header', 'quote', 'codeBlock', 'code', 'link', 'undo', 'redo']}
                        {...props}
                        value={value}
                        ref={ref}
                        enableScroll={false}
                        style={{
                            fontFamily: 'Posterama regular'
                        }}
                    />
                </Tab>
                <Tab eventKey="preview" title="Preview" className='py-1'>
                    {
                        value ? (
                            <MarkdownEditor.Markdown  {...props} source={value} className="bg-white text-black ps-1 pt-1 fs-5" />

                        ) : <span className="fst-italic fs-6">Nothing to preview</span>
                    }
                </Tab>
            </Tabs>
        </>
    )
}

export default React.forwardRef(TextEditor)