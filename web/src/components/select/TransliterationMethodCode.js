import React from 'react';

const TransliterationMethodCode = () => {
  return (
    <>
    <option value={0}>Arabic (Arabic language)</option>
    <option value={1}>Arabic (Persian language)</option>
    <option value={2}>Armenian</option>
    <option value={3}>Cyrillic</option>
    <option value={4}>Devanagari & related Indic</option>
    <option value={5}>Georgian</option>
    <option value={6}>Greek</option>
    <option value={7}>Han (Hanzi, Kanji, Hanja)</option>
    <option value={8}>Hebrew</option>
    <option value={10}>Kana</option>
    <option value={11}>Korean</option>
    <option value={12}>Thai</option>
    <option value={13}>Unspecified Standard</option>
    </>
  )
}

export default TransliterationMethodCode;