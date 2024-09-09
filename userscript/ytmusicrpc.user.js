// ==UserScript==
// @name         YoutubeMusicDiscordRichPresence2024SuperPuperShitCode
// @match        *://music.youtube.com/*
// @grant        GM.xmlHttpRequest
// @icon         https://www.google.com/s2/favicons?sz=16&domain=music.youtube.com
// @icon64       https://www.google.com/s2/favicons?sz=64&domain=music.youtube.com
// @downloadURL  https://github.com/oddyamill/ytmusicrpc/raw/master/userscript/ytmusicrpc.user.js
// @updateURL    https://github.com/oddyamill/ytmusicrpc/raw/master/userscript/ytmusicrpc.user.js
// ==/UserScript==

(async function () {
  const authKey =
    'e8ab39d4b23d2877af508538de8424fd7c8ea4734870f462591b759acdf07199'

  function requestRpc(method, body) {
    const headers = {
      Authorization: authKey,
    }

    let data

    if (body !== undefined) {
      data = JSON.stringify(body)
      headers['Content-Type'] = 'application/json'
    }

    return GM.xmlHttpRequest({
      url: 'http://localhost:32484/rpc',
      method,
      headers,
      data,
    })
  }

  function updatePresence(track) {
    return requestRpc('POST', track)
  }

  function deletePresence() {
    return requestRpc('DELETE')
  }

  function sleep(timeout) {
    return new Promise((resolve) => setTimeout(resolve, timeout))
  }

  function parse(time) {
    const [sec, min, hour] = time
      .split(':')
      .map((t) => +t)
      .reverse()
    return (sec + min * 60 + (hour || 0) * 3600) * 1000
  }

  await sleep(1000)

  const video = document.querySelector('video.video-stream')

  function listener() {
    if (video.paused) {
      return deletePresence()
    }

    const [current, end] = document
      .querySelector('#left-controls > span')
      .textContent.trim()
      .split(' / ')

    const track = {
      trackId: document
        .querySelector('a.ytp-title-link.yt-uix-sessionlink')
        .href.match(/v=([^&#]{5,})/)?.[1],
      title:
        document.querySelector('.ytmusic-player-bar > .title')?.title ||
        navigator.mediaSession.metadata.title,
      artist: navigator.mediaSession.metadata.artist,
      artwork: navigator.mediaSession.metadata.artwork.at(-1)?.src,
      album:
        navigator.mediaSession.metadata.album ||
        [...document.querySelectorAll('.byline a')].at(-1)?.textContent ||
        undefined,
      current: null,
      end: null,
    }

    if (end !== undefined) {
      track.current = parse(current)
      track.end = parse(end)
    }

    updatePresence(track)
  }

  video.addEventListener('playing', listener)
  video.addEventListener('pause', listener)

  if (video.src) {
    listener()
  }
})()
