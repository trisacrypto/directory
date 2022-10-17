/* eslint-disable no-console */
/* eslint-disable no-eq-null */
/* eslint-disable new-cap */
/* eslint-disable global-require */
const PO = require('pofile');
const path = require('path');

const LOCAL = 'en-dh';
const PO_FILE_PATH = path.join(__dirname, `src/locales/${LOCAL}/messages.po`);

function trim(string) {
  return string.replace(/^\s+|\s+$/g, '');
}

function replatEmptyMsgStrByDshes(items) {
  return items.map(item => ({
    ...item,
    msgstr: item.msgstr.toString() ? item.msgstr : ['---']
  }
  ));
}

function save(po) {
  require('fs').writeFile(PO_FILE_PATH, po, 'utf8', (error) => {
    if (error) {
      console.error('[saving]', error);
    }

    console.log(`file saved successfully in /locales/${LOCAL}/messages.po`);
  });
}

/*
 serialize to po
*/

// reverse what extract(string) method during PO.parse does
function _escape(str) {
  // don't unescape \n, since string can never contain it
  // since split('\n') is called on it
  // eslint-disable-next-line no-control-regex
  const string = str.replace(/[\x07\b\t\v\f\r"\\]/g, (match) => {
    switch (match) {
      case '\x07':
        return '\\a';
      case '\b':
        return '\\b';
      case '\t':
        return '\\t';
      case '\v':
        return '\\v';
      case '\f':
        return '\\f';
      case '\r':
        return '\\r';
      default:
        return '\\' + match;
    }
  });
  return string;
}

function _process(keyword, text, i) {
  // eslint-disable-next-line @typescript-eslint/no-shadow
  const lines = [];
  const parts = text.split(/\n/);
  const index = typeof i !== 'undefined' ? '[' + i + ']' : '';
  if (parts.length > 1) {
    lines.push(keyword + index + ' ""');
    parts.forEach((part) => {
      lines.push('"' + _escape(part) + '"');
    });
  } else {
    lines.push(keyword + index + ' "' + _escape(text) + '"');
  }
  return lines;
}

function _processLineBreak(keyword, text, index) {
  const processed = _process(keyword, text, index);
  for (let i = 1; i < processed.length - 1; i++) {
    processed[i] = processed[i].slice(0, -1) + '\\n"';
  }
  return processed;
}

function serializeToPo(options) {
  let lines = [];
  const self = options;

  // handle \n in single-line texts (can not be handled in _escape)

  // https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
  // says order is translator-comments, extracted-comments, references, flags

  options.comments.forEach((c) => {
    lines.push('# ' + c);
  });

  options.extractedComments.forEach((c) => {
    lines.push('#. ' + c);
  });

  options.references.forEach((ref) => {
    lines.push('#: ' + ref);
  });

  const flags = Object.keys(options.flags).filter((flag) => {
    return !!options.flags[flag];
  }, options);
  if (flags.length > 0) {
    lines.push('#, ' + flags.join(','));
  }
  const mkObsolete = options.obsolete ? '#~ ' : '';

  ['msgctxt', 'msgid', 'msgid_plural', 'msgstr'].forEach((keyword) => {
    let text = self[keyword];
    // eslint-disable-next-line eqeqeq
    if (text != null) {
      let hasTranslation = false;
      if (Array.isArray(text)) {
        // eslint-disable-next-line @typescript-eslint/no-shadow
        hasTranslation = text.some((text) => {
          return text;
        });
      }

      if (Array.isArray(text) && text.length > 1) {
        text.forEach((t, i) => {
          const processed = _processLineBreak(keyword, t, i);
          lines = lines.concat(mkObsolete + processed.join('\n' + mkObsolete));
        });
      } else if (self.msgid_plural && keyword === 'msgstr' && !hasTranslation) {
        for (let pluralIndex = 0; pluralIndex < self.nplurals; pluralIndex++) {
          lines = lines.concat(mkObsolete + _process(keyword, '', pluralIndex));
        }
      } else {
        const index = (self.msgid_plural && Array.isArray(text)) ? 0 : undefined;
        text = Array.isArray(text) ? text.join() : text;
        const processed = _processLineBreak(keyword, text, index);
        lines = lines.concat(mkObsolete + processed.join('\n' + mkObsolete));
      }
    }
  });

  return lines.join('\n');
}


function parsePo(options) {
  const lines = [];

  // -------------

  if (options.comments) {
    options.comments.forEach((comment) => {
      trim(lines.push(('# ' + comment)));
    });
  }

  if (options.extractedComments) {
    options.extractedComments.forEach((comment) => {
      trim(lines.push(('#. ' + comment)));
    });
  }

  lines.push('msgid ""');
  lines.push('msgstr ""');

  // -------------

  const self = options;
  const headerOrder = [];

  let keys = [];
  if (options.headers) {
    keys = Object.keys(options.headers);
  }

  keys.forEach((key) => {
    if (headerOrder.indexOf(key) === -1) {
      headerOrder.push(key);
    }
  });

  headerOrder.forEach((key) => {
    lines.push('"' + key + ': ' + self.headers[key] + '\\n"');
  });

  lines.push('');

  options.items.forEach((item) => {
    lines.push(serializeToPo(item));
    lines.push('');
  });

  return lines.join('\n');
}

PO.load(PO_FILE_PATH, (err, po) => {
  if (err) {
    console.error('[po loading error]', err);
  }

  po.items = [...replatEmptyMsgStrByDshes(po.items)];

  save(parsePo(po));
});
