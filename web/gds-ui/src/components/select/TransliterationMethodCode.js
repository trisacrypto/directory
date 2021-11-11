import React from 'react';
import { t } from "@lingui/macro";


const TransliterationMethodCode = () => {
  
  return (
    <>
    <option value={1}>{t`Arabic (Arabic language)`}</option>
    <option value={2}>{t`Arabic (Persian language)`}</option>
    <option value={3}>{t`Armenian`}</option>
    <option value={4}>{t`Cyrillic`}</option>
    <option value={5}>{t`Devanagari & related Indic`}</option>
    <option value={6}>{t`Georgian`}</option>
    <option value={7}>{t`Greek`}</option>
    <option value={8}>{t`Han (Hanzi, Kanji, Hanja)`}</option>
    <option value={9}>{t`Hebrew`}</option>
    <option value={10}>{t`Kana`}</option>
    <option value={11}>{t`Korean`}</option>
    <option value={12}>{t`Thai`}</option>
    <option value={0}>{t`Unspecified Standard`}</option>
    </>
  )
}

export default TransliterationMethodCode;