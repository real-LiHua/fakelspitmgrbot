// https://lsposed.jasonkhew96.dev/js/main.js
;(() => {
  /**
   * @typedef {object} ValidateResponse
   * @property {boolean} ok
   * @property {string} [message]
   * @property {string} [challenge_code]
   */

  /**
   * @typedef {object} SubmitResponse
   * @property {boolean} ok
   * @property {string} [message]
   */

  /**
   * @typedef {object} BottomButton
   * @property {string} type
   * @property {string} text
   * @property {string} color
   * @property {string} textColor
   * @property {boolean} isVisible
   * @property {boolean} isActive
   * @property {boolean} hasShineEffect
   * @property {string} position
   * @property {boolean} isProgressVisible
   * @property {(text: string) => BottomButton} setText
   * @property {(callback: () => void) => BottomButton} onClick
   * @property {(callback: () => void) => BottomButton} offClick
   * @property {() => BottomButton} show
   * @property {() => BottomButton} hide
   * @property {() => BottomButton} enable
   * @property {() => BottomButton} disable
   * @property {(leaveAction: boolean) => BottomButton} showProgress
   * @property {() => BottomButton} hideProgress
   */

  /**
   * @typedef {object} WebAppObj
   * @property {string} initData
   * @property {string} version
   * @property {string} platform
   * @property {BottomButton} MainButton
   * @property {BottomButton} SecondaryButton
   * @property {(string) => boolean} isVersionAtLeast
   * @property {(message: string, callback: () => void) => void} showAlert
   * @property {(message: string, callback: (ok: boolean) => void) => void} showConfirm
   * @property {() => void} ready
   * @property {() => void} expand
   * @property {() => void} close
   */

  /**
   * @typedef {object} TelegramObj
   * @property {WebAppObj} WebApp
   */

  /**
   * @typedef {object} CustomWindowObj
   * @property {TelegramObj} Telegram
   */

  /**
   * @typedef {Window & CustomWindowObj} CustomWindow
   */

  const genRandStr = l => {
    let result = ''
    const characters =
      'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    const charactersLength = characters.length
    let counter = 0
    while (counter < l) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength))
      counter += 1
    }
    return result
  }

  /** @type CustomWindow */
  const w = window
  const webapp = w.Telegram.WebApp

  const checkUsername = (/** @type string */ username) => {
    return /^(?:[a-z\d](?:[a-z\d]|-(?=[a-z\d])){0,38}|[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*(_[a-zA-Z0-9]+))$/i.test(
      username
    )
  }

  const checkSignature = (/** @type string */ signature) => {
    return /^-----BEGIN SSH SIGNATURE-----[\s\S]*-----END SSH SIGNATURE-----$/.test(
      signature
    )
  }

  const getUsername = () => {
    return document.getElementById('username').value.trim()
  }

  const getSignature = () => {
    return document.getElementById('signature').value.trim()
  }

  const enableInput = (/** @type {boolean} */ enabled) => {
    const username = document.getElementById('username')
    const signature = document.getElementById('signature')
    username.disabled = !enabled
    signature.disabled = !enabled
  }

  const onSubmitted = (/** @type {SubmitResponse} */ data) => {
    if (data.ok) {
      showAlert("Approved!\n申请通过!", () => {
        webapp.close()
      })
      return
    }
    showAlert(data.message)
    setTimeout(() => {
      webapp.close()
    }, 10000)
  }

  const onSubmit = () => {
    hideAlert()
    const username = getUsername()
    const signature = getSignature()
    if (!checkUsername(username)) {
      showAlert('Invalid username\n用户名无效')
      return
    }
    if (!checkSignature(signature)) {
      showAlert('Invalid signature\n签名无效')
      return
    }
    webapp.MainButton.disable().showProgress()
    const body = {
      username: username,
      signature: signature
    }
    fetch('/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Auth': webapp.initData
      },
      body: JSON.stringify(body)
    })
      .then(resp => {
        if (
          resp.ok &&
          resp.headers.get('Content-Type') === 'application/json'
        ) {
          return resp.json()
        }
        throw new Error('Failed to submit')
      })
      .then(data => {
        onSubmitted(data)
      })
      .catch(err => {
        showAlert(err.message)
      })
  }

  const setupMainButton = () => {
    webapp.MainButton.setText('Submit')
      .onClick(() => {
        webapp.showConfirm('Are you sure?\n确定?', ok => {
          if (ok) {
            onSubmit()
          }
        })
      })
      .hide()
      .showProgress(false)
      .disable()
  }

  const displayLoading = (/** @type boolean */ display) => {
    const overlay = document.getElementById('loading-overlay')
    if (display) {
      overlay.classList.remove('hide')
      return
    }
    overlay.classList.add('hide')
  }

  const showAlertCompat = (
    /** @type string */ message,
    /** @type {() => void} */ callback
  ) => {
    /** @type HTMLElement */
    const el = document.getElementsByClassName('alert').item(0)
    el.innerHTML = ''

    const splits = message.split('\n')
    for (const split of splits) {
      const p = document.createElement('p')
      p.innerText = split
      el.appendChild(p)
    }

    el.classList.add('show')
    callback ? setTimeout(() => {
      callback()
    }, 10000) : null
  }

  const showAlert = (
    /** @type {string} */ msg,
    /** @type {() => void} */ callback
  ) => {
    if (!webapp.isVersionAtLeast('6.2')) {
      showAlertCompat(msg, callback)
      return
    }
    webapp.showAlert(msg, callback)
  }

  const hideAlert = () => {
    const el = document.getElementsByClassName('alert').item(0)
    el.classList.remove('show')
  }

  const onAuthValid = (/** @type ValidateResponse */ data) => {
    ;[...document.getElementsByClassName('challenge-code')].forEach(
      /** @type HTMLSpanElement */ el => {
        el.innerText = data.challenge_code
        el.classList.add('unblur-challenge-code')
      }
    )
    displayLoading(false)
    enableInput(true)
    setupCopy()
    webapp.MainButton.enable().hideProgress().show()
    webapp.expand()
  }

  const onAuthInvalid = (/** @type ValidateResponse */ data) => {
    showAlert(data.message, () => {
      webapp.close()
    })
    displayLoading(false)
    webapp.MainButton.hide()
  }

  const onValidateResponse = (/** @type ValidateResponse */ data) => {
    if (data.ok) {
      onAuthValid(data)
    } else {
      onAuthInvalid(data)
    }
  }

  const onLoad = () => {
    webapp.ready()
    setupMainButton()

    fetch('/validate', {
      headers: {
        'X-Auth': webapp.initData
      }
    })
      .then(resp => {
        if (
          resp.ok &&
          resp.headers.get('Content-Type') === 'application/json'
        ) {
          return resp.json()
        }
        throw new Error('Failed to validate\n验证失败')
      })
      .then((/** @type ValidateResponse */ data) => {
        onValidateResponse(data)
      })
      .catch(err => {
        displayLoading(false)
        showAlert(err.message, () => {
          webapp.close()
        })
      })
  }

  const setupCopy = () => {
    ;[...document.getElementsByClassName('command-copy')].forEach(el => {
      el.addEventListener('click', () => {
        navigator.clipboard.writeText(el.innerText.trim())
      })
    })
  }

  onLoad()
})()

