function getQuery (params) {
  let newParams = new Map()
  for (let [k, v] of Object.entries(params)) {
    if (v != null && v != '') {
      newParams.set(k, v)
    }
  }

  return new URLSearchParams(newParams).toString()
}

function getFullPathWithQuery (path, params) {
  let query = getQuery(params)

  let fullPath = path
  if (query != '') {
    fullPath = fullPath + '?' + query
  }

  return fullPath
}

function cleanQuery (params) {
  let newParams = new Map()
  for (let [k, v] of Object.entries(params)) {
    if (v != null && v != '') {
      newParams.set(k, v)
    }
  }

  return Object.fromEntries(newParams)
}

export { getFullPathWithQuery, cleanQuery }
