<?php
/**
 * 工具类
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 12:02
 */
namespace app\Service;
use \zhimiao\Config;
use zhimiao\Data;
use Respect\Validation\Validator as v;
class Utils {

    /**
     * 缓存锁 - 此方法需要配合计划任务解决死锁问题
     * @param string $key 键名
     * @param bool $begin true-上锁 false-解锁
     * @param int $expire 超时时间
     * @return bool 伪布尔
     */
    public static function cacheNumLock($key = '', $num = 5, $scale_time = 1800, $lock_time = 60)
    {
        $redis = Data::redis();
        $key = 'cacheNumLock:'. $key;
        // 如果计次为0，则清空锁
        if ($num === 0) {
            return (bool) $redis->del($key);
        }
        $data = $redis->get($key);
        $data = @json_decode($data, true);
        if (empty($data)) {
            $data = [
                'num' => 1,
                'last_time' => 0,
            ];
        }
        $out_num = $data['num'] - $num;
        if ($out_num > 0) {
            // 锁定时间次方增长，直至超过超时时间
            $lock_time = pow($lock_time, $out_num);
            if (time() < $data['last_time'] + $lock_time) {
                return false;
            }
            $scale_time = max([$scale_time, $lock_time]);
        }
        $data['num']++;
        $data['last_time'] = time();
        return (bool) $redis->set($key, json_encode($data), $scale_time);
    }

    /**
     * 创建加密的sessionKey
     * @param $uid
     * @param int $expire
     * @return string
     */
    public static function createSessionKey($uid, $expire = 86400)
    {
        $data = json_encode([
            'uid' => $uid,
            'expire' => time() + $expire
        ]);
        return \zhimiao\Utils::encrypt($data);
    }

    /**
     * 解析sessionKey
     * @param $sessionKey
     * @param bool $isExpire
     * @return mixed|string
     */
    public static function parseSessionKey($sessionKey, $isExpire = true)
    {
        $data = \zhimiao\Utils::decrypt($sessionKey);
        if (!v::json()->validate($data)) {
            return 'sessionKey无法解析';
        }
        $data = json_decode($data, true);
        if ($isExpire && time() > $data['expire']) {
            return 'sessionKey已经过期了';
        }
        return $data;
    }
}