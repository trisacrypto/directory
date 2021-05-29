import React from 'react';

const TransliterationMethodCode = () => {
  return (
    <>
    <option value={1}>Arabic (Arabic language)</option>
    <option value={2}>Arabic (Persian language)</option>
    <option value={3}>Armenian</option>
    <option value={4}>Cyrillic</option>
    <option value={5}>Devanagari & related Indic</option>
    <option value={6}>Georgian</option>
    <option value={7}>Greek</option>
    <option value={8}>Han (Hanzi, Kanji, Hanja)</option>
    <option value={9}>Hebrew</option>
    <option value={10}>Kana</option>
    <option value={11}>Korean</option>
    <option value={12}>Thai</option>
    <option value={0}>Unspecified Standard</option>
    </>
  )
}

export default TransliterationMethodCode;