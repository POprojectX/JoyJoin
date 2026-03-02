<template>
  <div class="page">
    <div class="shell">
      <aside class="sidebar">
        <div class="logo" id="tgt-logo">
          <div class="logo-icon">
            <svg
              width="18"
              height="18"
              fill="none"
              stroke="white"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <circle cx="12" cy="12" r="10" />
              <path d="M8 14s1.5 2 4 2 4-2 4-2" />
              <line x1="9" y1="9" x2="9.01" y2="9" />
              <line x1="15" y1="9" x2="15.01" y2="9" />
            </svg>
          </div>
          <span class="logo-name">JoyJoin</span>
        </div>

        <div class="profile-card" id="tgt-profile">
          <div class="avatar-ring">
            <span class="avatar-letter">{{ user.name[0] }}</span>
          </div>
          <p class="profile-name">{{ user.name }}</p>
          <p class="profile-role">{{ user.role }}</p>
          <div class="profile-stats">
            <div class="pstat">
              <span class="pstat-val">14</span
              ><span class="pstat-lbl">Joined</span>
            </div>
            <div class="pstat-div" />
            <div class="pstat">
              <span class="pstat-val">3</span
              ><span class="pstat-lbl">Hosted</span>
            </div>
            <div class="pstat-div" />
            <div class="pstat">
              <span class="pstat-val">89</span
              ><span class="pstat-lbl">Friends</span>
            </div>
          </div>
          <button
            class="btn-create"
            id="tgt-create"
            @click="showCreateModal = true"
          >
            <svg
              width="13"
              height="13"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <path d="M12 5v14M5 12h14" />
            </svg>
            Create event
          </button>
        </div>

        <nav class="nav" id="tgt-nav">
          <span class="nav-label">Main</span>
          <button
            v-for="item in navItems"
            :key="item.id"
            class="nav-item"
            :class="{ active: activeTab === item.id }"
            @click="activeTab = item.id"
          >
            <svg
              width="18"
              height="18"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
              v-html="item.icon"
            />
            {{ item.label }}
            <span v-if="item.badge" class="nav-badge">{{ item.badge }}</span>
          </button>

          <div class="nav-divider" />
          <span class="nav-label">Account</span>
          <button class="nav-item" @click="activeTab = 'settings'">
            <svg
              width="18"
              height="18"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
            >
              <circle cx="12" cy="12" r="3" />
              <path
                d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"
              />
            </svg>
            Settings
          </button>
          <button class="nav-item nav-logout">
            <svg
              width="18"
              height="18"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
            >
              <path
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
            Sign out
          </button>
        </nav>

        <div class="sidebar-promo" id="tgt-promo">
          <div class="promo-glow" />
          <p class="promo-label">✨ Pro Plan</p>
          <p class="promo-text">
            Unlock unlimited events & priority discovery.
          </p>
          <button class="btn-upgrade">Upgrade now →</button>
        </div>
      </aside>

      <main class="main">
        <div class="hero" id="tgt-hero">
          <div class="hero-canvas">
            <div class="hero-blobs">
              <div class="blob blob1" />
              <div class="blob blob2" />
              <div class="blob blob3" />
              <div class="blob blob4" />
            </div>
            <div class="hero-grain" />
          </div>

          <div class="hero-top">
            <div class="search-wrap" id="tgt-search">
              <svg
                width="15"
                height="15"
                fill="none"
                stroke="white"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.35-4.35" />
              </svg>
              <input
                v-model="searchQuery"
                placeholder="Search events, venues, people…"
                class="search-input"
              />
            </div>
            <button
              class="notif-btn"
              id="tgt-notif"
              @click="notifOpen = !notifOpen"
            >
              <svg
                width="18"
                height="18"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 01-3.46 0" />
              </svg>
              <span class="notif-dot" />
            </button>
          </div>

          <div class="hero-bottom">
            <p class="hero-label">{{ todayStr }}</p>
            <h1 class="hero-greeting">
              Good {{ timeOfDay }}, <em>{{ user.name.split(" ")[0] }}</em> 👋
            </h1>
            <p class="hero-sub">Here's what's happening around you today.</p>
            <div class="hero-meta">
              <span class="hero-stat"
                ><span class="hero-stat-dot" style="background: #a78bfa" />{{
                  events.length
                }}
                upcoming events</span
              >
              <span class="hero-sep">·</span>
              <span class="hero-stat"
                ><span class="hero-stat-dot" style="background: #34d399" />{{
                  newsFeed.length
                }}
                new articles</span
              >
              <span class="hero-sep">·</span>
              <span class="hero-stat"
                ><span
                  class="hero-stat-dot"
                  style="background: #f59e0b"
                />Opole, PL</span
              >
            </div>
          </div>

          <div class="hero-ticker">
            <div class="ticker-track">
              <span v-for="t in tickerItems" :key="t" class="ticker-item">{{
                t
              }}</span>
              <span
                v-for="t in tickerItems"
                :key="'dup-' + t"
                class="ticker-item"
                aria-hidden="true"
                >{{ t }}</span
              >
            </div>
          </div>
        </div>

        <div class="filter-row" id="tgt-filters">
          <button
            v-for="f in filters"
            :key="f.id"
            class="filter-btn"
            :class="{ active: activeFilter === f.id }"
            @click="activeFilter = f.id"
          >
            {{ f.label }}
          </button>
        </div>

        <div class="content">
          <section class="col-news">
            <div class="section-head">
              <div>
                <h2 class="section-title">Latest News</h2>
                <p class="section-sub">
                  {{ filteredNews.length }} articles today
                </p>
              </div>
              <button class="btn-view-all">
                View all
                <svg
                  width="12"
                  height="12"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                >
                  <path d="m9 18 6-6-6-6" />
                </svg>
              </button>
            </div>
            <div class="news-list" id="tgt-news">
              <div
                v-for="(a, i) in filteredNews"
                :key="a.id"
                class="news-card"
                :style="{ animationDelay: i * 0.07 + 's' }"
              >
                <div class="news-thumb" :style="{ background: a.thumbBg }">
                  {{ a.emoji }}
                </div>
                <div class="news-body">
                  <div class="news-head-row">
                    <span
                      class="news-tag"
                      :style="{ background: a.tagBg, color: a.tagColor }"
                      >{{ a.tag }}</span
                    >
                    <span class="news-time">{{ a.time }}</span>
                  </div>
                  <p class="news-title">{{ a.title }}</p>
                  <p class="news-desc">{{ a.desc }}</p>
                  <div class="news-footer">
                    <div class="news-author">
                      <div
                        class="author-av"
                        :style="{ background: a.authorColor }"
                      >
                        {{ a.author[0] }}
                      </div>
                      <span>{{ a.author }}</span>
                    </div>
                    <button class="btn-read">Read →</button>
                  </div>
                </div>
              </div>
              <div v-if="filteredNews.length === 0" class="empty">
                No articles found
              </div>
            </div>
          </section>

          <section class="col-events">
            <div class="section-head">
              <div>
                <h2 class="section-title">Near You</h2>
                <p class="section-sub">Based on your location · Opole</p>
              </div>
              <button class="btn-view-all">
                View all
                <svg
                  width="12"
                  height="12"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                >
                  <path d="m9 18 6-6-6-6" />
                </svg>
              </button>
            </div>
            <div class="events-grid" id="tgt-events">
              <div
                v-for="(ev, i) in filteredEvents"
                :key="ev.id"
                class="event-card"
                :style="{ animationDelay: i * 0.08 + 's' }"
                @click="selectedEvent = ev"
              >
                <div class="event-img-area" :class="ev.bgClass">
                  <div
                    class="event-img-bg"
                    :style="{ background: ev.gradient }"
                  />
                  <span class="event-emoji">{{ ev.emoji }}</span>
                  <span
                    class="event-badge"
                    :class="{
                      tonight: ev.tag === 'Tonight',
                      soon: ev.tag === 'Tomorrow',
                    }"
                    >{{ ev.tag }}</span
                  >
                  <button
                    class="ev-fav"
                    :class="{ active: ev.favorited }"
                    @click.stop="ev.favorited = !ev.favorited"
                  >
                    <svg
                      width="12"
                      height="12"
                      :fill="ev.favorited ? 'currentColor' : 'none'"
                      stroke="currentColor"
                      stroke-width="2.5"
                      viewBox="0 0 24 24"
                    >
                      <path
                        d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"
                      />
                    </svg>
                  </button>
                </div>
                <div class="event-body">
                  <p class="event-title">{{ ev.title }}</p>
                  <p class="event-loc">
                    <svg
                      width="10"
                      height="10"
                      fill="var(--muted)"
                      viewBox="0 0 24 24"
                    >
                      <path
                        d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
                      />
                    </svg>
                    {{ ev.location }}
                  </p>
                  <div class="event-footer">
                    <div style="display: flex; align-items: center">
                      <div class="avatars">
                        <div
                          v-for="(c, n) in ev.colors"
                          :key="n"
                          class="mini-av"
                          :style="{ background: c }"
                        >
                          {{ "ABC"[n] }}
                        </div>
                      </div>
                      <span class="more-count">+{{ ev.going }}</span>
                    </div>
                    <button class="ev-join" @click.stop>Join</button>
                  </div>
                </div>
              </div>
              <div
                v-if="filteredEvents.length === 0"
                class="empty"
                style="grid-column: 1/-1"
              >
                No events found
              </div>
            </div>
          </section>

          <div class="col-full">
            <div class="featured-card" id="tgt-featured">
              <div class="featured-visual">
                <div class="featured-orbs">
                  <div class="forb forb1" />
                  <div class="forb forb2" />
                </div>
                <span class="featured-emoji">🎆</span>
              </div>
              <div class="featured-body">
                <div class="featured-eyebrow">
                  <span class="featured-tag">Featured · This weekend</span>
                  <span class="featured-live"
                    ><span class="live-dot" />Selling fast</span
                  >
                </div>
                <h3 class="featured-title">Summer Beach Party — Opole</h3>
                <p class="featured-desc">
                  Join hundreds of people at Proszkowski Beach for an
                  unforgettable night of music, fireworks, and good vibes. Open
                  bar until midnight, live DJ sets, and a laser show at 23:00.
                </p>
                <div class="featured-chips">
                  <span class="chip">🎵 Live Music</span>
                  <span class="chip">🎆 Fireworks</span>
                  <span class="chip">🍹 Open Bar</span>
                  <span class="chip">📍 Proszkowski Beach</span>
                </div>
                <div class="featured-footer">
                  <button class="btn-join">
                    <svg
                      width="14"
                      height="14"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.5"
                      viewBox="0 0 24 24"
                    >
                      <path d="M5 12h14M12 5l7 7-7 7" />
                    </svg>
                    Join Event
                  </button>
                  <div class="featured-social">
                    <div class="feat-avs">
                      <div
                        v-for="(c, i) in [
                          '#7C3AED',
                          '#A78BFA',
                          '#6D28D9',
                          '#8B5CF6',
                        ]"
                        :key="i"
                        class="mini-av"
                        :style="{ background: c }"
                      >
                        {{ "ABCD"[i] }}
                      </div>
                    </div>
                    <span class="featured-going">125+ people going</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="col-full stats-row" id="tgt-stats">
            <div v-for="s in quickStats" :key="s.label" class="stat-card">
              <div class="stat-icon">{{ s.icon }}</div>
              <div>
                <p class="stat-val">{{ s.value }}</p>
                <p class="stat-lbl">{{ s.label }}</p>
              </div>
              <div class="stat-trend" :class="s.trend > 0 ? 'up' : 'down'">
                {{ s.trend > 0 ? "↑" : "↓" }} {{ Math.abs(s.trend) }}%
              </div>
            </div>
          </div>

          <div class="col-full">
            <div class="cal-section" id="tgt-calendar">
              <div class="cal-header">
                <div>
                  <h2 class="cal-title">{{ calMonthName }} {{ calYear }}</h2>
                  <p class="cal-sub">{{ eventDays.size }} events this month</p>
                </div>
                <div class="cal-nav">
                  <button class="cal-nav-btn" @click="prevMonth">
                    <svg
                      width="13"
                      height="13"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.5"
                      viewBox="0 0 24 24"
                    >
                      <path d="m15 18-6-6 6-6" />
                    </svg>
                  </button>
                  <button class="cal-nav-btn" @click="nextMonth">
                    <svg
                      width="13"
                      height="13"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.5"
                      viewBox="0 0 24 24"
                    >
                      <path d="m9 18 6-6-6-6" />
                    </svg>
                  </button>
                </div>
              </div>
              <div class="cal-days-header">
                <div
                  v-for="d in ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']"
                  :key="d"
                  class="cal-day-name"
                >
                  {{ d }}
                </div>
              </div>
              <div class="cal-grid">
                <div
                  v-for="cell in calCells"
                  :key="cell.key"
                  class="cal-day"
                  :class="{
                    other: !cell.current,
                    today: cell.isToday,
                    'has-event': cell.hasEvent,
                    selected: cell.isSelected,
                  }"
                  @click="cell.current && (selectedCalDay = cell.day)"
                >
                  {{ cell.day }}
                  <span v-if="cell.hasEvent && !cell.isToday" class="ev-dot" />
                </div>
              </div>
              <div class="cal-legend">
                <span class="leg-item"
                  ><span
                    class="leg-dot"
                    style="background: var(--violet)"
                  />Today</span
                >
                <span class="leg-item"
                  ><span
                    class="leg-dot"
                    style="
                      background: var(--violet-m);
                      border: 1px solid var(--violet);
                    "
                  />Has event</span
                >
                <span class="leg-item"
                  ><span
                    class="leg-dot"
                    style="
                      background: var(--violet-s);
                      border: 1px solid var(--violet);
                    "
                  />Selected</span
                >
              </div>
            </div>
          </div>

          <div class="col-full friends-section">
            <div class="section-head">
              <div>
                <h2 class="section-title">Friends Going</h2>
                <p class="section-sub">See what your network is attending</p>
              </div>
              <button class="btn-view-all">
                Invite friends
                <svg
                  width="12"
                  height="12"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                >
                  <path d="M12 5v14M5 12h14" />
                </svg>
              </button>
            </div>
            <div class="friends-list" id="tgt-friends">
              <div v-for="f in friends" :key="f.name" class="friend-card">
                <div class="friend-av" :style="{ background: f.color }">
                  {{ f.name[0] }}
                </div>
                <div class="friend-info">
                  <p class="friend-name">{{ f.name }}</p>
                  <p class="friend-ev">{{ f.event }}</p>
                </div>
                <span class="friend-time">{{ f.time }}</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Transition name="modal">
      <div
        v-if="selectedEvent"
        class="modal-overlay"
        @click.self="selectedEvent = null"
      >
        <div class="modal-card">
          <div class="modal-hero" :class="selectedEvent.bgClass">
            <div
              class="modal-hero-bg"
              :style="{ background: selectedEvent.gradient }"
            />
            <span class="modal-emoji">{{ selectedEvent.emoji }}</span>
            <button class="modal-close" @click="selectedEvent = null">
              <svg
                width="16"
                height="16"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
                viewBox="0 0 24 24"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="modal-tags">
              <span
                class="event-badge"
                :class="{ tonight: selectedEvent.tag === 'Tonight' }"
                >{{ selectedEvent.tag }}</span
              >
              <span class="modal-cat">{{ selectedEvent.category }}</span>
            </div>
            <h3 class="modal-title">{{ selectedEvent.title }}</h3>
            <p class="modal-location">
              <svg
                width="14"
                height="14"
                fill="var(--muted)"
                viewBox="0 0 24 24"
              >
                <path
                  d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
                />
              </svg>
              {{ selectedEvent.location }}
            </p>
            <p class="modal-desc">{{ selectedEvent.desc }}</p>
            <div class="modal-footer">
              <button class="btn-join" style="flex: 1">
                <svg
                  width="14"
                  height="14"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                >
                  <path d="M5 12h14M12 5l7 7-7 7" />
                </svg>
                Join Event
              </button>
              <button class="btn-share">
                <svg
                  width="16"
                  height="16"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <circle cx="18" cy="5" r="3" />
                  <circle cx="6" cy="12" r="3" />
                  <circle cx="18" cy="19" r="3" />
                  <line x1="8.59" y1="13.51" x2="15.42" y2="17.49" />
                  <line x1="15.41" y1="6.51" x2="8.59" y2="10.49" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <div
      v-if="tutActive"
      class="tut-blocker"
      @wheel.prevent
      @touchmove.prevent
      @scroll.prevent
    />

    <Transition name="tut-fade">
      <div v-if="tutActive" class="tut-overlay-wrap">
        <svg class="tut-svg" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <mask id="tut-mask">
              <rect width="100%" height="100%" fill="white" />
              <rect
                :x="holeRect.x"
                :y="holeRect.y"
                :width="holeRect.w"
                :height="holeRect.h"
                :rx="holeRect.rx"
                fill="black"
              />
            </mask>
          </defs>
          <rect
            width="100%"
            height="100%"
            fill="rgba(12,0,30,0.76)"
            mask="url(#tut-mask)"
          />
        </svg>

        <div
          class="tut-pulse"
          :style="{
            left: holeRect.x - 6 + 'px',
            top: holeRect.y - 6 + 'px',
            width: holeRect.w + 12 + 'px',
            height: holeRect.h + 12 + 'px',
          }"
        />
      </div>
    </Transition>

    <Transition name="tut-card">
      <div
        v-if="tutActive"
        class="tut-card"
        :style="{ top: cardPos.top + 'px', left: cardPos.left + 'px' }"
      >
        <div class="tut-header">
          <div class="tut-step-badge">
            <span class="tut-step-dot" />
            Step {{ currentStep + 1 }} of {{ steps.length }}
          </div>
          <div class="tut-title">{{ steps[currentStep].title }}</div>
        </div>
        <div class="tut-body">
          <p class="tut-desc">{{ steps[currentStep].desc }}</p>

          <div class="tut-progress">
            <div
              v-for="i in steps.length"
              :key="i"
              class="tut-dot"
              :class="{
                active: i - 1 === currentStep,
                done: i - 1 < currentStep,
              }"
            />
          </div>

          <div class="tut-actions">
            <button class="tut-skip" @click="closeTutorial">
              Skip tutorial
            </button>
            <div class="tut-btns">
              <button
                v-if="currentStep > 0"
                class="tut-btn-prev"
                @click="prevStep"
              >
                ← Back
              </button>
              <button class="tut-btn-next" @click="nextStep">
                <template v-if="currentStep < steps.length - 1">
                  Next
                  <svg
                    width="12"
                    height="12"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                    viewBox="0 0 24 24"
                  >
                    <path d="M5 12h14M12 5l7 7-7 7" />
                  </svg>
                </template>
                <template v-else>🎉 Finish</template>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="tut-done">
      <div v-if="tutDone" class="tut-done-overlay">
        <div class="tut-done-card">
          <span class="tut-done-emoji">🎉</span>
          <div class="tut-done-title">You're all set!</div>
          <p class="tut-done-sub">
            You've completed the JoyJoin tour. Start exploring events near you
            and connect with friends right now.
          </p>
          <button class="tut-done-btn" @click="closeTutorial">
            <svg
              width="14"
              height="14"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <path d="M5 12h14M12 5l7 7-7 7" />
            </svg>
            Start exploring
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from "vue";

const activeTab = ref("hub");
const activeFilter = ref("all");
const searchQuery = ref("");
const calDate = ref(new Date());
const selectedEvent = ref(null);
const selectedCalDay = ref(null);
const notifOpen = ref(false);
const showCreateModal = ref(false);

const user = { name: "Pashok Krutoi", role: "Event Organizer" };

const timeOfDay = computed(() => {
  const h = new Date().getHours();
  if (h < 12) return "morning";
  if (h < 18) return "afternoon";
  return "evening";
});

const todayStr = computed(() =>
  new Date().toLocaleDateString("en-US", {
    weekday: "long",
    month: "long",
    day: "numeric",
    year: "numeric",
  }),
);

const navItems = [
  {
    id: "hub",
    label: "My Hub",
    icon: '<path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>',
  },
  {
    id: "events",
    label: "Events",
    icon: '<rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>',
  },
  {
    id: "guests",
    label: "Guests",
    icon: '<path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/>',
  },
  {
    id: "messages",
    label: "Messages",
    icon: '<path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/>',
    badge: 3,
  },
];

const filters = [
  { id: "all", label: "All" },
  { id: "today", label: "Today" },
  { id: "weekend", label: "This Weekend" },
  { id: "music", label: "🎵 Music" },
  { id: "food", label: "🍴 Food" },
  { id: "sports", label: "⚽ Sports" },
  { id: "art", label: "🎨 Art" },
  { id: "outdoor", label: "🌿 Outdoor" },
];

const tickerItems = [
  "Beach Bonfire Night tonight →",
  "125 going to Summer Beach Party →",
  "Food Festival returns to Opole →",
  "City Run 10K sign-ups open →",
  "Art & Wine Evening this Friday →",
];

const newsFeed = ref([
  {
    id: 1,
    emoji: "🎉",
    tag: "Venues",
    tagBg: "#ede9fe",
    tagColor: "#7c3aed",
    thumbBg: "linear-gradient(135deg,#ede9fe,#ddd6fe)",
    title: "Top 5 Party Venues in Warsaw",
    desc: "Discover the hottest event spaces before summer ends. Rooftops, galleries & more await.",
    time: "2h ago",
    author: "Marta Kowalska",
    authorColor: "#7c3aed",
  },
  {
    id: 2,
    emoji: "🍕",
    tag: "Food",
    tagBg: "#fef3c7",
    tagColor: "#b45309",
    thumbBg: "linear-gradient(135deg,#fef9c3,#fde68a)",
    title: "Food Festival Returns to Opole",
    desc: "80+ vendors, live cooking battles and music. The biggest edition in festival history.",
    time: "5h ago",
    author: "Tomasz Wiśniewski",
    authorColor: "#f59e0b",
  },
  {
    id: 3,
    emoji: "🎸",
    tag: "Music",
    tagBg: "#ecfdf5",
    tagColor: "#065f46",
    thumbBg: "linear-gradient(135deg,#d1fae5,#a7f3d0)",
    title: "Live Concert Series Announced",
    desc: "Three-day outdoor festival with local bands and international headliners.",
    time: "1d ago",
    author: "Anna Nowak",
    authorColor: "#10b981",
  },
  {
    id: 4,
    emoji: "🏃",
    tag: "Sports",
    tagBg: "#fff7ed",
    tagColor: "#c2410c",
    thumbBg: "linear-gradient(135deg,#ffedd5,#fed7aa)",
    title: "City Run 10K — Sign Ups Open",
    desc: "The annual city run is back. Register early to secure your spot.",
    time: "1d ago",
    author: "Kamil Piotrowiak",
    authorColor: "#ef4444",
  },
  {
    id: 5,
    emoji: "🎨",
    tag: "Art",
    tagBg: "#fdf4ff",
    tagColor: "#7e22ce",
    thumbBg: "linear-gradient(135deg,#fae8ff,#f3e8ff)",
    title: "Gallery Night: Local Artists Showcase",
    desc: "Seven emerging artists open their studios to the public. Wine, live painting.",
    time: "2d ago",
    author: "Zofia Jabłońska",
    authorColor: "#a855f7",
  },
]);

const events = ref([
  {
    id: 1,
    emoji: "🏖️",
    title: "Beach Bonfire Night",
    location: "Proszkowski Beach",
    tag: "Tonight",
    going: 42,
    bgClass: "ev-bg-a",
    gradient: "linear-gradient(135deg,#4c1d95,#7c3aed)",
    colors: ["#7C3AED", "#A78BFA", "#6B7280"],
    category: "Outdoor",
    desc: "Gather around the fire, enjoy live acoustic sets, roasted marshmallows and a stunning night by the water.",
    favorited: false,
    tags: ["outdoor", "today"],
  },
  {
    id: 2,
    emoji: "🎨",
    title: "Art & Wine Evening",
    location: "Old Town Gallery",
    tag: "Fri 7 PM",
    going: 18,
    bgClass: "ev-bg-b",
    gradient: "linear-gradient(135deg,#1e1b4b,#4338ca)",
    colors: ["#6D28D9", "#8B5CF6", "#9CA3AF"],
    category: "Art",
    desc: "Paint, sip and socialise. A curated evening pairing local wine with guided canvas painting sessions.",
    favorited: false,
    tags: ["art", "weekend"],
  },
  {
    id: 3,
    emoji: "🎮",
    title: "LAN Gaming Party",
    location: "CyberHub, Opole",
    tag: "Sat",
    going: 67,
    bgClass: "ev-bg-c",
    gradient: "linear-gradient(135deg,#1a1a2e,#5b21b6)",
    colors: ["#5B21B6", "#7C3AED", "#6B7280"],
    category: "Gaming",
    desc: "Bring your rig or rent one. 12 hours of competitive and casual gaming across all genres.",
    favorited: false,
    tags: ["weekend"],
  },
  {
    id: 4,
    emoji: "🌿",
    title: "Park Yoga Morning",
    location: "Botanical Garden",
    tag: "Sun 9 AM",
    going: 24,
    bgClass: "ev-bg-d",
    gradient: "linear-gradient(135deg,#064e3b,#065f46)",
    colors: ["#8B5CF6", "#C4B5FD", "#6B7280"],
    category: "Outdoor",
    desc: "Start your Sunday with a restorative flow session surrounded by nature. All levels welcome.",
    favorited: false,
    tags: ["outdoor", "weekend"],
  },
  {
    id: 5,
    emoji: "🎵",
    title: "Jazz Under the Stars",
    location: "City Amphitheatre",
    tag: "Tomorrow",
    going: 98,
    bgClass: "ev-bg-e",
    gradient: "linear-gradient(135deg,#1c1917,#44403c)",
    colors: ["#92400e", "#d97706", "#6B7280"],
    category: "Music",
    desc: "An open-air jazz concert featuring three ensembles performing classic standards.",
    favorited: false,
    tags: ["music", "today"],
  },
  {
    id: 6,
    emoji: "🍜",
    title: "Asian Street Food Night",
    location: "Market Square",
    tag: "Sat 5 PM",
    going: 134,
    bgClass: "ev-bg-f",
    gradient: "linear-gradient(135deg,#7f1d1d,#b91c1c)",
    colors: ["#dc2626", "#ef4444", "#6B7280"],
    category: "Food",
    desc: "Explore flavours from Thailand, Japan, Vietnam and Korea. 20+ food stalls.",
    favorited: false,
    tags: ["food", "weekend"],
  },
]);

const quickStats = [
  { icon: "🎟️", label: "Events Joined", value: "14", trend: 12 },
  { icon: "👥", label: "New Connections", value: "28", trend: 8 },
  { icon: "⭐", label: "Your Rating", value: "4.9", trend: 2 },
  { icon: "🏆", label: "Achievements", value: "7", trend: 0 },
];

const friends = [
  {
    name: "Marta Kowalska",
    event: "Beach Bonfire Night",
    time: "Tonight",
    color: "#7c3aed",
  },
  {
    name: "Tomek Wiśniewski",
    event: "LAN Gaming Party",
    time: "Saturday",
    color: "#10b981",
  },
  {
    name: "Anna Nowak",
    event: "Jazz Under the Stars",
    time: "Tomorrow",
    color: "#f59e0b",
  },
  {
    name: "Kamil Piotrowiak",
    event: "City Run 10K",
    time: "Sunday",
    color: "#ef4444",
  },
  {
    name: "Zofia Jabłońska",
    event: "Art & Wine Evening",
    time: "Friday",
    color: "#8b5cf6",
  },
];

const eventDays = new Set([1, 5, 8, 12, 15, 19, 22, 26, 28]);

const filteredNews = computed(() => {
  const q = searchQuery.value.toLowerCase();
  return q
    ? newsFeed.value.filter(
        (n) =>
          n.title.toLowerCase().includes(q) || n.desc.toLowerCase().includes(q),
      )
    : newsFeed.value;
});

const filteredEvents = computed(() => {
  const q = searchQuery.value.toLowerCase();
  const f = activeFilter.value;
  let list = events.value;
  if (f !== "all")
    list = list.filter(
      (e) => e.tags.includes(f) || e.category.toLowerCase() === f,
    );
  if (q)
    list = list.filter(
      (e) =>
        e.title.toLowerCase().includes(q) ||
        e.location.toLowerCase().includes(q),
    );
  return list;
});

const monthNames = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];
const calMonth = computed(() => calDate.value.getMonth());
const calYear = computed(() => calDate.value.getFullYear());
const calMonthName = computed(() => monthNames[calMonth.value]);

const calCells = computed(() => {
  const y = calYear.value,
    m = calMonth.value;
  const firstDay = new Date(y, m, 1).getDay();
  const daysInMonth = new Date(y, m + 1, 0).getDate();
  const daysInPrev = new Date(y, m, 0).getDate();
  const today = new Date();
  const cells = [];
  for (let i = firstDay - 1; i >= 0; i--)
    cells.push({
      key: `p${i}`,
      day: daysInPrev - i,
      current: false,
      isToday: false,
      hasEvent: false,
      isSelected: false,
    });
  for (let d = 1; d <= daysInMonth; d++) {
    const isToday =
      d === today.getDate() &&
      m === today.getMonth() &&
      y === today.getFullYear();
    cells.push({
      key: `c${d}`,
      day: d,
      current: true,
      isToday,
      hasEvent: eventDays.has(d),
      isSelected: selectedCalDay.value === d,
    });
  }
  const rem = 42 - cells.length;
  for (let d = 1; d <= rem; d++)
    cells.push({
      key: `n${d}`,
      day: d,
      current: false,
      isToday: false,
      hasEvent: false,
      isSelected: false,
    });
  return cells;
});

function prevMonth() {
  const d = new Date(calDate.value);
  d.setMonth(d.getMonth() - 1);
  calDate.value = d;
}
function nextMonth() {
  const d = new Date(calDate.value);
  d.setMonth(d.getMonth() + 1);
  calDate.value = d;
}

const tutActive = ref(false);
const tutDone = ref(false);
const currentStep = ref(0);

const holeRect = ref({ x: 0, y: 0, w: 0, h: 0, rx: 14 });
const cardPos = ref({ top: 0, left: 0 });

const steps = [
  {
    target: "tgt-logo",
    title: "Welcome to JoyJoin 🎉",
    desc: "JoyJoin is your all-in-one social hub for discovering, joining, and hosting local events. This quick tour walks you through every feature in under 2 minutes.",
    pad: 12,
    side: "right",
  },
  {
    target: "tgt-profile",
    title: "Your Profile Card",
    desc: "See your personal stats at a glance — events joined, events hosted, and the size of your friends network. Click any stat to drill in further.",
    pad: 8,
    side: "right",
  },
  {
    target: "tgt-create",
    title: "Create an Event",
    desc: 'Ready to host? Hit "Create event" to set up your own gathering — add a description, invite friends, pick a location, and go live instantly.',
    pad: 6,
    side: "right",
  },
  {
    target: "tgt-nav",
    title: "Main Navigation",
    desc: "Switch between your Hub, the Events calendar, Guests manager, and Messages from here. The badge on Messages shows unread conversations.",
    pad: 8,
    side: "right",
  },
  {
    target: "tgt-search",
    title: "Search Everything",
    desc: "Type anything — event names, venues, or people — to instantly filter both the news feed and the events grid in real time.",
    pad: 10,
    side: "bottom",
  },
  {
    target: "tgt-filters",
    title: "Filter by Category",
    desc: "Narrow down what you see by time (Today, This Weekend) or by vibe (Music, Food, Sports, Art, Outdoor). Combine with search for laser-sharp results.",
    pad: 8,
    side: "bottom",
  },
  {
    target: "tgt-news",
    title: "Latest News Feed",
    desc: "Stay in the loop with curated event news from local authors. Each article is tagged by category and links directly to the full story.",
    pad: 10,
    side: "right",
  },
  {
    target: "tgt-events",
    title: "Events Near You",
    desc: 'Browse events close to your location. Click any card for full details. Hit the ♡ heart to save favourites, or "Join" to RSVP straight away.',
    pad: 10,
    side: "left",
  },
  {
    target: "tgt-featured",
    title: "Featured Event",
    desc: "The spotlight card highlights the must-attend event of the week. Chips show exactly what's included. Spots fill fast — grab yours early!",
    pad: 12,
    side: "top",
  },
  {
    target: "tgt-stats",
    title: "Your Activity Stats",
    desc: "Track how active you've been — events attended, new friends made, your community rating, and achievements unlocked.",
    pad: 10,
    side: "top",
  },
  {
    target: "tgt-calendar",
    title: "Event Calendar",
    desc: "Navigate months to see upcoming events. Days with a dot have events scheduled. Click a day to see what's on. Today is always highlighted in violet.",
    pad: 10,
    side: "top",
  },
  {
    target: "tgt-friends",
    title: "Friends Going",
    desc: "See which friends are attending events this week. Staying connected is effortless — if a friend joins, you're more likely to show up too.",
    pad: 10,
    side: "top",
  },
  {
    target: "tgt-promo",
    title: "✨ Go Pro",
    desc: "Upgrade to Pro for unlimited event creation, priority discovery ranking, advanced analytics, and exclusive community badges.",
    pad: 8,
    side: "right",
  },
];

function lockPage() {
  document.body.style.userSelect = "none";
}
function unlockPage() {
  document.body.style.overflow = "";
  document.body.style.userSelect = "";
}

function getTargetRect(id) {
  return document.getElementById(id)?.getBoundingClientRect() ?? null;
}

function positionCard(rect, side) {
  const cardEl = document.querySelector(".tut-card");
  const cardW = cardEl?.offsetWidth || 340;
  const cardH = cardEl?.offsetHeight || 260;
  const gap = 22;
  const vw = window.innerWidth;
  const vh = window.innerHeight;
  let top, left;

  if (side === "right") {
    top = rect.top + rect.height / 2 - cardH / 2;
    left = rect.right + gap;
    if (left + cardW > vw - 12) {
      left = rect.left - cardW - gap;
    }
  } else if (side === "left") {
    top = rect.top + rect.height / 2 - cardH / 2;
    left = rect.left - cardW - gap;
    if (left < 12) {
      left = rect.right + gap;
    }
  } else if (side === "bottom") {
    top = rect.bottom + gap;
    left = rect.left + rect.width / 2 - cardW / 2;
  } else {
    // top
    top = rect.top - cardH - gap;
    left = rect.left + rect.width / 2 - cardW / 2;
  }

  top = Math.max(12, Math.min(top, vh - cardH - 12));
  left = Math.max(12, Math.min(left, vw - cardW - 12));
  return { top, left };
}

async function renderStep(idx) {
  const step = steps[idx];
  if (!step) return;

  const el = document.getElementById(step.target);
  if (el) el.scrollIntoView({ behavior: "smooth", block: "nearest" });

  await new Promise((r) => setTimeout(r, 180));

  const rect = getTargetRect(step.target);
  if (!rect) return;

  const pad = step.pad ?? 10;

  holeRect.value = {
    x: rect.left - pad,
    y: rect.top - pad,
    w: rect.width + pad * 2,
    h: rect.height + pad * 2,
    rx: 14,
  };

  await nextTick();
  cardPos.value = positionCard(rect, step.side);
}

async function startTutorial() {
  currentStep.value = 0;
  tutActive.value = true;
  tutDone.value = false;
  lockPage();
  await nextTick();
  renderStep(0);
}

async function nextStep() {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++;
    await renderStep(currentStep.value);
  } else {
    // Finish
    tutActive.value = false;
    tutDone.value = true;
  }
}

async function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--;
    await renderStep(currentStep.value);
  }
}

function closeTutorial() {
  tutActive.value = false;
  tutDone.value = false;
  unlockPage();
}

function onKeyDown(e) {
  if (!tutActive.value && !tutDone.value) return;
  if (tutDone.value) {
    if (e.key === "Escape" || e.key === "Enter") closeTutorial();
    return;
  }
  if (e.key === "ArrowRight" || e.key === "Enter") nextStep();
  if (e.key === "ArrowLeft") prevStep();
  if (e.key === "Escape") closeTutorial();
}

function onResize() {
  if (tutActive.value) renderStep(currentStep.value);
}

onMounted(() => {
  document.addEventListener("keydown", onKeyDown);
  window.addEventListener("resize", onResize);
  // Auto-start after a short delay
  setTimeout(startTutorial, 700);
});

onBeforeUnmount(() => {
  document.removeEventListener("keydown", onKeyDown);
  window.removeEventListener("resize", onResize);
  unlockPage();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=DM+Serif+Display:ital@0;1&family=DM+Sans:opsz,wght@9..40,300;9..40,400;9..40,500;9..40,600;9..40,700&display=swap");

:root {
  --bg: #f5f3ff;
  --bg2: #ede9fe;
  --surface: #ffffff;
  --border: #e5e0f8;
  --violet: #7c3aed;
  --violet2: #6d28d9;
  --violet3: #5b21b6;
  --violet-s: #ede9fe;
  --violet-m: #ddd6fe;
  --gold: #a78bfa;
  --teal: #8b5cf6;
  --ink: #1e1433;
  --ink2: #3b2f68;
  --muted: #6e5fa8;
  --soft: #a898d4;
  --sh-sm: 0 2px 12px rgba(109, 40, 217, 0.08);
  --sh-md: 0 8px 28px rgba(109, 40, 217, 0.13);
  --sh-lg: 0 20px 56px rgba(109, 40, 217, 0.18);
  --r-xl: 24px;
  --r-lg: 18px;
  --r-md: 12px;
  --r-sm: 8px;
  --font-d: "DM Serif Display", serif;
  --font-b: "DM Sans", sans-serif;
}

*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
button {
  font-family: var(--font-b);
  cursor: pointer;
  border: none;
  background: none;
}
input {
  font-family: var(--font-b);
}

.page {
  min-height: 100vh;
  background:
    radial-gradient(ellipse at 0% 0%, var(--bg2) 0%, transparent 45%),
    radial-gradient(ellipse at 100% 0%, var(--border) 0%, transparent 45%),
    var(--bg);
  display: flex;
  align-items: stretch;
  font-family: var(--font-b);
  color: var(--ink);
}
.shell {
  display: flex;
  align-items: stretch;
  width: 100%;
  min-height: 100vh;
}

.sidebar {
  width: 268px;
  flex-shrink: 0;
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  padding: 28px 18px 20px;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
}
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
}
.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--violet), var(--violet3));
  border-radius: var(--r-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(109, 40, 217, 0.3);
}
.logo-name {
  font-family: var(--font-d);
  font-size: 22px;
  color: var(--ink);
}

.profile-card {
  background: var(--bg);
  border-radius: var(--r-lg);
  padding: 18px;
  text-align: center;
  margin-bottom: 24px;
  border: 1px solid var(--border);
}
.avatar-ring {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--violet), #a78bfa);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 10px;
  box-shadow:
    0 0 0 3px var(--surface),
    0 0 0 5px var(--border);
}
.avatar-letter {
  font-family: var(--font-d);
  font-size: 22px;
  color: white;
}
.profile-name {
  font-size: 14px;
  font-weight: 600;
}
.profile-role {
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}
.profile-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin: 12px 0 0;
  padding: 10px 0 0;
  border-top: 1px solid var(--border);
}
.pstat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}
.pstat-val {
  font-size: 15px;
  font-weight: 700;
  line-height: 1;
}
.pstat-lbl {
  font-size: 10px;
  color: var(--soft);
}
.pstat-div {
  width: 1px;
  height: 24px;
  background: var(--border);
}

.btn-create {
  margin-top: 12px;
  width: 100%;
  background: linear-gradient(135deg, var(--violet), var(--violet3));
  color: white;
  border-radius: var(--r-md);
  padding: 9px 16px;
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  box-shadow: 0 4px 16px rgba(109, 40, 217, 0.3);
  transition:
    transform 0.15s,
    box-shadow 0.2s;
}
.btn-create:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 22px rgba(109, 40, 217, 0.4);
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}
.nav-label {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--soft);
  padding: 8px 10px 4px;
  margin-top: 8px;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 11px 12px;
  border-radius: var(--r-md);
  font-size: 14px;
  font-weight: 500;
  color: var(--muted);
  transition: all 0.18s;
}
.nav-item svg {
  flex-shrink: 0;
  opacity: 0.7;
  transition: opacity 0.18s;
}
.nav-item:hover,
.nav-item.active {
  background: var(--violet-s);
  color: var(--violet);
}
.nav-item:hover svg,
.nav-item.active svg {
  opacity: 1;
}
.nav-badge {
  margin-left: auto;
  background: var(--violet);
  color: white;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 20px;
}
.nav-divider {
  height: 1px;
  background: var(--border);
  margin: 10px 0;
}
.nav-logout {
  color: var(--soft);
}
.nav-logout:hover {
  background: #fdf2ff;
  color: #9333ea;
}

.sidebar-promo {
  margin-top: 12px;
  background: linear-gradient(135deg, #2e1065, #4c1d95);
  border-radius: var(--r-lg);
  padding: 16px;
  position: relative;
}
.promo-glow {
  align-items: center;
  justify-content: center;
  display: flex;
  inset: 0;
  background: radial-gradient(
    circle at 80% 20%,
    rgba(167, 139, 250, 0.3),
    transparent 60%
  );
}
.promo-label {
  font-size: 11px;
  font-weight: 700;
  color: #c4b5fd;
  margin-bottom: 4px;
  position: relative;
}
.promo-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
  line-height: 1.4;
  position: relative;
}
.btn-upgrade {
  margin-top: 10px;
  width: 100%;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  border-radius: var(--r-sm);
  padding: 8px 14px;
  font-size: 12px;
  font-weight: 600;
  position: relative;
  transition: background 0.2s;
}
.btn-upgrade:hover {
  background: rgba(255, 255, 255, 0.2);
}

.main {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.hero {
  position: relative;
  display: flex;
  flex-direction: column;
}
.hero-canvas {
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, #2e1065 0%, #4c1d95 45%, #3b0764 100%);
}
.hero-blobs {
  position: absolute;
  inset: 0;
}
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.4;
  animation: blobFloat var(--dur, 8s) ease-in-out infinite alternate;
}
.blob1 {
  width: 320px;
  height: 320px;
  background: #7c3aed;
  top: -80px;
  right: 15%;
  --dur: 7s;
}
.blob2 {
  width: 250px;
  height: 250px;
  background: #a78bfa;
  top: -40px;
  right: 42%;
  --dur: 9s;
}
.blob3 {
  width: 200px;
  height: 200px;
  background: #5b21b6;
  bottom: -60px;
  right: 5%;
  --dur: 6s;
}
.blob4 {
  width: 180px;
  height: 180px;
  background: #c4b5fd;
  top: 20px;
  left: 5%;
  --dur: 11s;
  opacity: 0.2;
}
.hero-grain {
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  opacity: 0.05;
}
@keyframes blobFloat {
  from {
    transform: translate(0, 0) scale(1);
  }
  to {
    transform: translate(18px, 14px) scale(1.06);
  }
}

.hero-top {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 36px 0;
}
.search-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 14px;
  padding: 10px 18px;
  transition:
    background 0.2s,
    border-color 0.2s;
}
.search-wrap:focus-within {
  background: rgba(255, 255, 255, 0.18);
  border-color: rgba(255, 255, 255, 0.3);
}
.search-input {
  background: transparent;
  border: none;
  outline: none;
  font-size: 14px;
  color: white;
  flex: 1;
}
.search-input::placeholder {
  color: rgba(255, 255, 255, 0.45);
}
.notif-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: background 0.2s;
}
.notif-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}
.notif-dot {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #f59e0b;
  border: 2px solid #4c1d95;
}

.hero-bottom {
  position: relative;
  z-index: 2;
  padding: 20px 36px 16px;
}
.hero-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 6px;
}
.hero-greeting {
  font-family: var(--font-d);
  font-size: 36px;
  color: white;
  line-height: 1.1;
  letter-spacing: -0.02em;
}
.hero-greeting em {
  color: #c4b5fd;
  font-style: normal;
}
.hero-sub {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.45);
  margin-top: 6px;
}
.hero-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
  flex-wrap: wrap;
}
.hero-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.65);
}
.hero-stat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.hero-sep {
  color: rgba(255, 255, 255, 0.25);
  font-size: 13px;
}

.hero-ticker {
  position: relative;
  z-index: 2;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(8px);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding: 8px 0;
}
.ticker-track {
  display: flex;
  gap: 40px;
  animation: ticker 30s linear infinite;
  white-space: nowrap;
}
.ticker-item {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
}
@keyframes ticker {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-50%);
  }
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px 36px 0;
  overflow-x: auto;
}
.filter-row::-webkit-scrollbar {
  display: none;
}
.filter-btn {
  padding: 7px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: var(--muted);
  background: var(--surface);
  border: 1px solid var(--border);
  white-space: nowrap;
  flex-shrink: 0;
  transition: all 0.18s;
}
.filter-btn:hover {
  border-color: var(--violet-m);
  color: var(--violet);
}
.filter-btn.active {
  background: var(--violet);
  color: white;
  border-color: var(--violet);
  font-weight: 600;
}

.content {
  padding: 28px 36px 40px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  flex: 1;
}
.col-news {
  grid-column: 1;
}
.col-events {
  grid-column: 2;
}
.col-full {
  grid-column: 1/-1;
}

.section-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 14px;
}
.section-title {
  font-family: var(--font-d);
  font-size: 20px;
  line-height: 1.1;
}
.section-sub {
  font-size: 12px;
  color: var(--soft);
  margin-top: 2px;
}
.btn-view-all {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--violet);
  padding: 5px 12px;
  border-radius: var(--r-sm);
  background: var(--violet-s);
  transition: background 0.18s;
}
.btn-view-all:hover {
  background: var(--violet-m);
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.news-card {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  background: var(--surface);
  border-radius: var(--r-lg);
  padding: 14px;
  border: 1px solid var(--border);
  box-shadow: var(--sh-sm);
  animation: fadeUp 0.45s ease both;
  transition:
    transform 0.22s,
    box-shadow 0.22s,
    border-color 0.22s;
  cursor: pointer;
}
.news-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--sh-md);
  border-color: var(--violet-m);
}
.news-thumb {
  width: 46px;
  height: 46px;
  border-radius: var(--r-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.news-body {
  flex: 1;
  min-width: 0;
}
.news-head-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
}
.news-tag {
  display: inline-block;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: 2px 8px;
  border-radius: 20px;
}
.news-time {
  font-size: 11px;
  color: var(--soft);
}
.news-title {
  font-size: 13px;
  font-weight: 600;
  line-height: 1.35;
}
.news-desc {
  font-size: 12px;
  color: var(--muted);
  margin-top: 3px;
  line-height: 1.5;
}
.news-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}
.news-author {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--muted);
}
.author-av {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: white;
}
.btn-read {
  font-size: 12px;
  font-weight: 600;
  color: var(--violet);
  padding: 4px 10px;
  border-radius: 20px;
  background: var(--violet-s);
}

.events-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.event-card {
  background: var(--surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--border);
  box-shadow: var(--sh-sm);
  animation: fadeUp 0.45s ease both;
  transition:
    transform 0.22s,
    box-shadow 0.22s;
  cursor: pointer;
}
.event-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--sh-md);
}
.event-img-area {
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
.event-img-bg {
  position: absolute;
  inset: 0;
  opacity: 0.85;
}
.event-emoji {
  font-size: 36px;
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.3));
}
.event-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(6px);
  color: white;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 20px;
}
.event-badge.tonight {
  background: var(--violet);
}
.event-badge.soon {
  background: #059669;
}
.ev-fav {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 2;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(4px);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition:
    background 0.2s,
    color 0.2s;
}
.ev-fav:hover {
  background: rgba(239, 68, 68, 0.7);
}
.ev-fav.active {
  background: rgba(239, 68, 68, 0.85);
  color: #fca5a5;
}
.event-body {
  padding: 10px 12px 12px;
}
.event-title {
  font-size: 13px;
  font-weight: 600;
  line-height: 1.3;
}
.event-loc {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: var(--muted);
  margin-top: 3px;
}
.event-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}
.avatars {
  display: flex;
}
.mini-av {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  font-size: 9px;
  font-weight: 700;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--surface);
  margin-right: -6px;
}
.more-count {
  font-size: 11px;
  font-weight: 600;
  color: var(--muted);
  margin-left: 10px;
}
.ev-join {
  font-size: 11px;
  font-weight: 600;
  color: var(--violet);
  background: var(--violet-s);
  padding: 3px 9px;
  border-radius: 20px;
  transition: background 0.18s;
}
.ev-join:hover {
  background: var(--violet-m);
}
.ev-bg-a {
  background: #f5f0ff;
}
.ev-bg-b {
  background: #eef2ff;
}
.ev-bg-c {
  background: #f0fdf4;
}
.ev-bg-d {
  background: #faf5ff;
}
.ev-bg-e {
  background: #fefce8;
}
.ev-bg-f {
  background: #fff1f2;
}

.featured-card {
  display: flex;
  align-items: stretch;
  background: linear-gradient(135deg, #2e1065, #4c1d95);
  border-radius: var(--r-xl);
  box-shadow: var(--sh-lg);
  min-height: 190px;
  transition:
    transform 0.22s,
    box-shadow 0.22s;
}
.featured-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 28px 64px rgba(109, 40, 217, 0.3);
}
.featured-visual {
  width: 210px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: linear-gradient(135deg, #1e0a4e, #3b0764);
}
.featured-orbs {
  position: absolute;
  inset: 0;
}
.forb {
  position: absolute;
  border-radius: 50%;
  filter: blur(30px);
  animation: blobFloat var(--dur, 6s) ease-in-out infinite alternate;
}
.forb1 {
  width: 120px;
  height: 120px;
  background: rgba(124, 58, 237, 0.6);
  top: -20px;
  left: -20px;
  --dur: 5s;
}
.forb2 {
  width: 100px;
  height: 100px;
  background: rgba(196, 181, 253, 0.4);
  bottom: -10px;
  right: -10px;
  --dur: 7s;
}
.featured-emoji {
  position: relative;
  z-index: 1;
  font-size: 72px;
  filter: drop-shadow(0 4px 20px rgba(0, 0, 0, 0.4));
}
.featured-body {
  flex: 1;
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.featured-eyebrow {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}
.featured-tag {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #c4b5fd;
  background: rgba(196, 181, 253, 0.12);
  padding: 3px 10px;
  border-radius: 20px;
  border: 1px solid rgba(196, 181, 253, 0.25);
}
.featured-live {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  color: #4ade80;
}
.live-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #4ade80;
  animation: livePulse 1.5s ease-in-out infinite;
}
@keyframes livePulse {
  0%,
  100% {
    opacity: 1;
    box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.4);
  }
  50% {
    opacity: 0.8;
    box-shadow: 0 0 0 5px rgba(74, 222, 128, 0);
  }
}
.featured-title {
  font-family: var(--font-d);
  font-size: 26px;
  color: white;
  line-height: 1.15;
  letter-spacing: -0.01em;
  margin-bottom: 8px;
}
.featured-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
  line-height: 1.6;
  max-width: 420px;
}
.featured-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 12px;
}
.chip {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.65);
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  padding: 3px 10px;
  border-radius: 20px;
}
.featured-footer {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 18px;
}
.btn-join {
  display: flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #7c3aed, #5b21b6);
  color: white;
  border-radius: var(--r-md);
  padding: 10px 22px;
  font-size: 13px;
  font-weight: 700;
  box-shadow: 0 6px 20px rgba(109, 40, 217, 0.4);
  transition:
    transform 0.15s,
    box-shadow 0.2s;
}
.btn-join:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 28px rgba(109, 40, 217, 0.5);
}
.featured-social {
  display: flex;
  align-items: center;
  gap: 8px;
}
.feat-avs {
  display: flex;
}
.featured-going {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}
.stat-card {
  background: var(--surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--border);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: var(--sh-sm);
  transition:
    transform 0.2s,
    box-shadow 0.2s;
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--sh-md);
}
.stat-icon {
  font-size: 24px;
}
.stat-val {
  font-size: 20px;
  font-weight: 700;
  font-family: var(--font-d);
  line-height: 1;
}
.stat-lbl {
  font-size: 11px;
  color: var(--soft);
  margin-top: 2px;
}
.stat-trend {
  margin-left: auto;
  font-size: 11px;
  font-weight: 700;
  padding: 3px 7px;
  border-radius: 20px;
}
.stat-trend.up {
  background: #ecfdf5;
  color: #059669;
}
.stat-trend.down {
  background: #fff1f2;
  color: #e11d48;
}

.cal-section {
  background: var(--surface);
  border-radius: var(--r-xl);
  border: 1px solid var(--border);
  padding: 24px;
}
.cal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16px;
}
.cal-title {
  font-family: var(--font-d);
  font-size: 20px;
}
.cal-sub {
  font-size: 12px;
  color: var(--soft);
  margin-top: 2px;
}
.cal-nav {
  display: flex;
  gap: 4px;
}
.cal-nav-btn {
  width: 30px;
  height: 30px;
  border-radius: var(--r-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
  color: var(--muted);
  transition:
    background 0.18s,
    color 0.18s;
}
.cal-nav-btn:hover {
  background: var(--violet-s);
  color: var(--violet);
}
.cal-days-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 6px;
}
.cal-day-name {
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  color: var(--soft);
  padding: 4px 0;
}
.cal-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}
.cal-day {
  aspect-ratio: 1;
  border-radius: var(--r-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--muted);
  position: relative;
  cursor: pointer;
  transition:
    background 0.15s,
    color 0.15s;
}
.cal-day:hover {
  background: var(--bg);
  color: var(--ink);
}
.cal-day.other {
  color: var(--border);
  cursor: default;
}
.cal-day.other:hover {
  background: transparent;
  color: var(--border);
}
.cal-day.today {
  background: var(--violet);
  color: white;
  font-weight: 700;
}
.cal-day.selected:not(.today) {
  background: var(--violet-s);
  color: var(--violet);
  font-weight: 600;
}
.ev-dot {
  position: absolute;
  bottom: 4px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--violet);
}
.cal-legend {
  display: flex;
  gap: 16px;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}
.leg-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--muted);
}
.leg-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.friends-list {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
}
.friend-card {
  background: var(--surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--border);
  padding: 14px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 6px;
  box-shadow: var(--sh-sm);
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  cursor: pointer;
}
.friend-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--sh-md);
}
.friend-av {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-d);
  font-size: 18px;
  color: white;
}
.friend-info {
  flex: 1;
  min-width: 0;
}
.friend-name {
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.friend-ev {
  font-size: 11px;
  color: var(--muted);
  margin-top: 2px;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.friend-time {
  font-size: 10px;
  font-weight: 600;
  color: var(--violet);
  background: var(--violet-s);
  padding: 2px 8px;
  border-radius: 20px;
}

.empty {
  text-align: center;
  padding: 28px;
  color: var(--soft);
  font-size: 14px;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(30, 20, 51, 0.7);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
}
.modal-card {
  background: var(--surface);
  border-radius: var(--r-xl);
  width: 100%;
  max-width: 440px;
  box-shadow: 0 40px 80px rgba(0, 0, 0, 0.4);
}
.modal-hero {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
.modal-hero-bg {
  position: absolute;
  inset: 0;
}
.modal-emoji {
  font-size: 72px;
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 4px 20px rgba(0, 0, 0, 0.3));
}
.modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 2;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}
.modal-close:hover {
  background: rgba(0, 0, 0, 0.6);
}
.modal-body {
  padding: 20px 24px 24px;
}
.modal-tags {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}
.modal-cat {
  font-size: 11px;
  font-weight: 600;
  color: var(--muted);
  background: var(--bg);
  padding: 2px 8px;
  border-radius: 20px;
}
.modal-title {
  font-family: var(--font-d);
  font-size: 22px;
  line-height: 1.2;
  margin-bottom: 6px;
}
.modal-location {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 10px;
}
.modal-desc {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.6;
}
.modal-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 18px;
}
.btn-share {
  width: 42px;
  height: 42px;
  border-radius: var(--r-md);
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  flex-shrink: 0;
  transition:
    border-color 0.2s,
    color 0.2s;
}
.btn-share:hover {
  border-color: var(--violet-m);
  color: var(--violet);
}

.modal-enter-active,
.modal-leave-active {
  transition:
    opacity 0.25s,
    transform 0.25s;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

@keyframes fadeUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tut-blocker {
  position: fixed;
  inset: 0;
  z-index: 7900;
  cursor: default;
  /* Prevent touch scroll on mobile */
  touch-action: none;
  overscroll-behavior: contain;
}

.tut-overlay-wrap {
  position: fixed;
  inset: 0;
  z-index: 8000;
  pointer-events: none; /* clicks pass through to tut-card / tut-blocker */
}
.tut-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  overflow: visible;
}

.tut-pulse {
  position: fixed;
  z-index: 8010;
  pointer-events: none;
  border-radius: 20px;
  animation: tutPulse 2.2s ease-in-out infinite;
}
@keyframes tutPulse {
  0%,
  100% {
    box-shadow: 0 0 0 0 rgba(124, 58, 237, 0.55);
  }
  50% {
    box-shadow: 0 0 0 14px rgba(124, 58, 237, 0);
  }
}

.tut-card {
  position: fixed;
  z-index: 8100;
  width: 340px;
  background: #ffffff;
  border-radius: 20px;
  box-shadow:
    0 32px 80px rgba(30, 0, 60, 0.28),
    0 0 0 1px rgba(124, 58, 237, 0.15);
}

.tut-header {
  background: linear-gradient(135deg, #7c3aed, #5b21b6);
  border-radius: 20px 20px 0 0;
  padding: 20px 22px 18px;
  position: relative;
}
.tut-header::before {
  content: "";
  position: absolute;
  inset: 0;
  background: radial-gradient(
    circle at 80% 30%,
    rgba(196, 181, 253, 0.35),
    transparent 60%
  );
}
.tut-step-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.25);
  color: white;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  padding: 3px 10px;
  border-radius: 20px;
  margin-bottom: 10px;
  position: relative;
}
.tut-step-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #c4b5fd;
  flex-shrink: 0;
}
.tut-title {
  font-family: var(--font-d);
  font-size: 21px;
  color: white;
  line-height: 1.2;
  position: relative;
}

.tut-body {
  padding: 18px 22px 22px;
}
.tut-desc {
  font-size: 13.5px;
  color: var(--muted);
  line-height: 1.7;
}

.tut-progress {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-top: 16px;
}
.tut-dot {
  height: 7px;
  border-radius: 4px;
  background: var(--border);
  transition:
    background 0.25s,
    width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 7px;
}
.tut-dot.active {
  background: var(--violet);
  width: 22px;
}
.tut-dot.done {
  background: var(--violet-m);
}

.tut-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
}
.tut-skip {
  font-size: 12px;
  font-weight: 500;
  color: var(--soft);
  padding: 6px 10px;
  border-radius: var(--r-sm);
  transition:
    color 0.18s,
    background 0.18s;
}
.tut-skip:hover {
  color: var(--muted);
  background: var(--bg);
}
.tut-btns {
  display: flex;
  gap: 8px;
}
.tut-btn-prev {
  padding: 9px 16px;
  border-radius: var(--r-md);
  font-size: 13px;
  font-weight: 600;
  color: var(--violet);
  background: var(--violet-s);
  transition: background 0.18s;
}
.tut-btn-prev:hover {
  background: var(--violet-m);
}
.tut-btn-next {
  padding: 9px 20px;
  border-radius: var(--r-md);
  font-size: 13px;
  font-weight: 700;
  color: white;
  background: linear-gradient(135deg, var(--violet), var(--violet3));
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 4px 16px rgba(109, 40, 217, 0.35);
  transition:
    transform 0.15s,
    box-shadow 0.2s;
}
.tut-btn-next:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 22px rgba(109, 40, 217, 0.45);
}

/* Completion screen */
.tut-done-overlay {
  position: fixed;
  inset: 0;
  z-index: 9000;
  background: rgba(12, 0, 30, 0.72);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
}
.tut-done-card {
  background: white;
  border-radius: 28px;
  padding: 48px 44px;
  text-align: center;
  max-width: 360px;
  box-shadow: 0 40px 80px rgba(30, 0, 60, 0.3);
}
.tut-done-emoji {
  font-size: 56px;
  display: block;
  margin-bottom: 16px;
}
.tut-done-title {
  font-family: var(--font-d);
  font-size: 26px;
  margin-bottom: 8px;
}
.tut-done-sub {
  font-size: 14px;
  color: var(--muted);
  line-height: 1.6;
}
.tut-done-btn {
  margin-top: 24px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, var(--violet), var(--violet3));
  color: white;
  border-radius: var(--r-md);
  padding: 12px 28px;
  font-size: 14px;
  font-weight: 700;
  font-family: var(--font-b);
  box-shadow: 0 6px 20px rgba(109, 40, 217, 0.35);
  cursor: pointer;
  border: none;
  transition:
    transform 0.15s,
    box-shadow 0.2s;
}
.tut-done-btn:hover {
  transform: translateY(-2px);
}

/* Tutorial transitions */
.tut-fade-enter-active,
.tut-fade-leave-active {
  transition: opacity 0.35s ease;
}
.tut-fade-enter-from,
.tut-fade-leave-to {
  opacity: 0;
}

.tut-card-enter-active {
  transition:
    opacity 0.35s ease,
    transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.tut-card-leave-active {
  transition: opacity 0.2s ease;
}
.tut-card-enter-from {
  opacity: 0;
  transform: scale(0.88) translateY(10px);
}
.tut-card-leave-to {
  opacity: 0;
}

.tut-done-enter-active {
  transition: opacity 0.35s ease;
}
.tut-done-leave-active {
  transition: opacity 0.25s ease;
}
.tut-done-enter-from,
.tut-done-leave-to {
  opacity: 0;
}

::-webkit-scrollbar {
  width: 5px;
  height: 5px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: var(--violet-m);
}

@media (max-width: 1100px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .friends-list {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (max-width: 900px) {
  .sidebar {
    display: none;
  }
  .content {
    grid-template-columns: 1fr;
    padding: 20px;
  }
  .col-news,
  .col-events {
    grid-column: 1;
  }
  .events-grid {
    grid-template-columns: 1fr 1fr;
  }
  .hero-top,
  .hero-bottom {
    padding-left: 20px;
    padding-right: 20px;
  }
  .filter-row {
    padding-left: 20px;
    padding-right: 20px;
  }
  .featured-visual {
    display: none;
  }
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .friends-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 560px) {
  .events-grid {
    grid-template-columns: 1fr;
  }
  .stats-row {
    grid-template-columns: 1fr 1fr;
  }
  .friends-list {
    grid-template-columns: 1fr 1fr;
  }
  .hero-greeting {
    font-size: 28px;
  }
  .tut-card {
    width: calc(100vw - 24px);
  }
}
</style>
