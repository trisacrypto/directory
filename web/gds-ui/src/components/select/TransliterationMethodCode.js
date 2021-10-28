import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const TransliterationMethodCode = () => {
  
  return (
    <>
    <option value={1}>{i18n._(t`Arabic (Arabic language)`)}</option>
    <option value={2}>{i18n._(t`Arabic (Persian language)`)}</option>
    <option value={3}>{i18n._(t`Armenian`)}</option>
    <option value={4}>{i18n._(t`Cyrillic`)}</option>
    <option value={5}>{i18n._(t`Devanagari & related Indic`)}</option>
    <option value={6}>{i18n._(t`Georgian`)}</option>
    <option value={7}>{i18n._(t`Greek`)}</option>
    <option value={8}>{i18n._(t`Han (Hanzi, Kanji, Hanja)`)}</option>
    <option value={9}>{i18n._(t`Hebrew`)}</option>
    <option value={10}>{i18n._(t`Kana`)}</option>
    <option value={11}>{i18n._(t`Korean`)}</option>
    <option value={12}>{i18n._(t`Thai`)}</option>
    <option value={0}>{i18n._(t`Unspecified Standard`)}</option>
    </>
  )
}

export default TransliterationMethodCode;